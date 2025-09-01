package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/genai"
)

func main() {
	appStart := time.Now()
	dbAddress := getEnv("DB_ADDRESS", "localhost:5432")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "")
	dbName := getEnv("DB_NAME", "svc-applications")
	dbArgs := getEnv("DB_ARGS", "sslmode=disable")
	proxiesEnv := getEnv("TRUSTED_PROXIES", "127.0.0.1")

	log.Printf("Подключение к postgres://%s:XXXXX@%s/%s?%s\n", dbUser, dbAddress, dbName, dbArgs)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?%s", dbUser, dbPass, dbAddress, dbName, dbArgs)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Ошибка создания подключения:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	log.Println("Подключение успешно! База данных работает.")

	log.Println("Проверка таблицы 'applications'")

	createTableSQL := `
		DO $$
		BEGIN
    		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'application_decision_enum') THEN
        		CREATE TYPE application_decision_enum AS ENUM ('Pending', 'Accept', 'Reject');
    		END IF;
		END
		$$;

		create table if not exists "applications" (
  		"id" SERIAL not null,
  		"nickname" varchar(16) not null check (char_length(nickname) between 3 and 16 and nickname ~ '^[a-zA-Z0-9_]{3,16}$'),
  		"age" smallint not null check (age BETWEEN 1 AND 80),
  		"about" varchar(4096) not null,
  		"join_reason" varchar(1024) not null,
  		"discord_nickname" varchar(32) not null check (discord_nickname ~ '^[\w\-]{2,32}$'),
  		"inviter" varchar(256) null,
  		"ai_decision" application_decision_enum not null default 'Pending',
  		"admin_decision" application_decision_enum not null default 'Pending',
  		"ai_answer" VARCHAR(1024),
  		"ai_comment" VARCHAR(4096),
  		constraint "users_pkey" primary key ("id")
		);

		comment on column "applications"."nickname" is 'Валидный никнейм minecraft';
		comment on column "applications"."age" is 'Возраст игрока 1-80 лет';
		comment on column "applications"."about" is 'Поле "О себе" длинной не более 4096 символов';
		comment on column "applications"."join_reason" is 'Поле "Почему хотите вступить" длинной не более 1024 символов';
		comment on column "applications"."discord_nickname" is 'Валидный discord ник, на старые ники типа example#1234 пофиг';
		comment on column "applications"."inviter" is 'Поле "Кто вас приласил" длинной не более 256 символов, опционально и может быть null';
		comment on column "applications"."ai_decision" is 'Мнение ИИ, принят или нет';
		comment on column "applications"."admin_decision" is 'Мнение админов принят или нет';
		comment on column "applications"."ai_answer" is 'Ответ ИИ самому игроку';
		comment on column "applications"."ai_comment" is 'Пояснение ИИ к своему ответу';
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Println("Ошибка при проверке таблицы 'applications': ", err)
		return
	}

	log.Println("Таблица 'applications' готова")

	log.Println("Проверка API Gemini")

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{ // API ключ берётся из ENV параметра GOOGLE_API_KEY
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
	requestStart := time.Now()

	result, err := client.Models.CountTokens(ctx, "gemini-2.5-flash-lite", genai.Text("Hello"), &genai.CountTokensConfig{})
	cancel()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Gemini работает, количество тестовых токенов: %d, пинг: %s\n", result.TotalTokens, time.Since(requestStart))

	log.Printf("Запуск Gin спустя: %s с начала запуска программы", time.Since(appStart))

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.SetTrustedProxies(strings.Split(proxiesEnv, ","))

	// Middleware для добавления db в контекст
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.Run(":9003")
}
