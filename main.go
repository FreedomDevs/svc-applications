package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "embed"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/genai"
)

//go:embed check_table.sql
var createTableSQL string

func main() {
	appStart := time.Now()
	dbAddress := getEnv("DB_ADDRESS", "localhost:8003")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "")
	dbName := getEnv("DB_NAME", "svc-applications")
	dbArgs := getEnv("DB_ARGS", "sslmode=disable")
	proxiesEnv := getEnv("TRUSTED_PROXIES", "127.0.0.1")
	listenPort := getEnv("LISTEN_PORT", "80")
	usersServiceUrl := getEnv("USERS_SERVICE_URL", "http://[fd98:2dd6:8f48:1d99:dc28:e6e1::2]:80")

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

	log.Println("Проверка API Gemini")

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{ // API ключ берётся из ENV параметра GOOGLE_API_KEY
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Подготовка gin")

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.SetTrustedProxies(strings.Split(proxiesEnv, ","))

	// Middleware для добавления db в контекст
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("gemini", client)
		c.Set("usersServiceUrl", usersServiceUrl)
		c.Next()
	})

	r.GET("/applications/", getAllHandler)
	r.GET("/applications/:uuid", getByUUIDHandler)
	r.GET("/applications/me", getMyHandler)
	r.POST("/applications/create", createApplicationHandler)

	log.Printf("Запуск Gin спустя: %s с начала запуска программы", time.Since(appStart))
	r.Run(":" + listenPort)
}
