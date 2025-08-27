package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	dbAddress := getEnv("DB_ADDRESS", "localhost:5432")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "")
	dbName := getEnv("DB_NAME", "svc-applications")
	dbArgs := getEnv("DB_ARGS", "sslmode=disable")

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?%s", dbUser, dbPass, dbAddress, dbName, dbArgs)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("Ошибка создания подключения:", err)
		return
	}
	defer db.Close()

	// Пробуем реально соединиться
	if err := db.Ping(); err != nil {
		fmt.Println("Ошибка подключения к БД:", err)
		return
	}

	fmt.Println("Подключение успешно! Данные валидны.")
}
