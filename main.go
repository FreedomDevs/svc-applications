package main

import (
	"context"
	"svc-applications/first_parser"
	//	"database/sql"
	_ "embed"
	//	"fmt"
	"log"
	//"strings"
	"time"

	//"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/genai"
)

//go:embed check_table.sql
var createTableSQL string

func main() {
	/*
		//appStart := time.Now()
		dbAddress := getEnv("DB_ADDRESS", "localhost:8003")
		dbUser := getEnv("DB_USER", "root")
		dbPass := getEnv("DB_PASS", "")
		dbName := getEnv("DB_NAME", "svc-applications")
		dbArgs := getEnv("DB_ARGS", "sslmode=disable")
		//proxiesEnv := getEnv("TRUSTED_PROXIES", "127.0.0.1")

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

		_, err = db.Exec(createTableSQL)
		if err != nil {
			fmt.Println("Ошибка при проверке таблицы 'applications': ", err)
			return
		}

		log.Println("Таблица 'applications' готова")
	*/

	log.Println("Проверка API Gemini")

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{ // API ключ берётся из ENV параметра GOOGLE_API_KEY
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	nickname := "mikinol"
	age := "16 лет"
	join_reason := "Ну я просто хочу поиграть на сервере подобно этому + у меня там играет друг PerchicYT."
	about := "Ну тип я Ярик, мне 16 лет. Я не особо ютубер, но если сервер понравится, может какие-то анимации буду делать. Я 2д аниматор. Я достаточно малообщителен, но люблю уделять время большим проектам.	Ну я просто хочу поиграть на сервере подобно этому + у меня там играет друг."
	inviter_by := "PerchicYT"

	result, err := first_parser.SendRequest(nickname, age, join_reason, about, inviter_by, ctx, client)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("%v", "\n"+result.Text())

	/*
		requestStart := time.Now()

			result, err := client.Models.CountTokens(ctx, "gemini-2.5-flash-lite", genai.Text("Hello"), &genai.CountTokensConfig{})
			cancel()
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Gemini работает, количество тестовых токенов: %d, пинг: %s\n", result.TotalTokens, time.Since(requestStart))



				log.Println("Подготовка gin")

				r := gin.New()

				r.Use(gin.Recovery())
				r.Use(gin.Logger())
				r.SetTrustedProxies(strings.Split(proxiesEnv, ","))

				// Middleware для добавления db в контекст
				r.Use(func(c *gin.Context) {
					c.Set("db", db)
					c.Set("gemini", client)
					c.Next()
				})

				route := r.Group("/applicatons")
				{
					_ = route
				}

				log.Printf("Запуск Gin спустя: %s с начала запуска программы", time.Since(appStart))
				r.Run(":9003")
	*/
}
