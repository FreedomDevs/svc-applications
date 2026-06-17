package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"svc-applications/final_reviewer"
	"svc-applications/first_parser"

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
	listenPort := getEnv("LISTEN_PORT", "9003")
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
		c.Next()
	})

	r.GET(":id", func(ctx *gin.Context) {
		db := ctx.MustGet("db").(*sql.DB)
		id := ctx.Param("id")

		var data ApplicationData
		err := db.QueryRow(`SELECT id,userid,age,about,join_reason,inviter,submitted_at,ai_categories,ai_decision,admin_decision,ai_answer,ai_comment
												FROM applications WHERE id = ?`, id).Scan(&data.id, &data.userid, &data.age, &data.about, &data.join_reason, &data.inviter,
			&data.submitted_at, &data.ai_categories, &data.ai_decision, &data.admin_decision, &data.ai_answer, &data.ai_comment)
		if err != nil {
			log.Panicln(err)
		}

		ctx.JSON(200, data)
	})

	r.POST("", func(ctx *gin.Context) {
		gemini := ctx.MustGet("gemini").(*genai.Client)
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("Connection", "keep-alive")
		ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

		id := ctx.GetHeader("eauth-user-id")
		if id == "" {
			ctx.JSON(422, "Пустой заголовок eauth-user-id")
			return
		}

		resp, err := http.Get(usersServiceUrl + "/users/" + id)
		if err != nil {
			log.Panicln(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Ошибка чтения:", err)
			return
		}

		var jsonBody map[string]any
		err = json.Unmarshal(body, &jsonBody)
		if err != nil {
			log.Println("Ошибка парсинга json:", err)
			return
		}

		username := jsonBody["data"].(map[string]any)["name"].(string)

		var req ApplicationRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(422, "Говно запрос, "+err.Error())
			return
		}

		ctx.Writer.Write([]byte("Классификация заявки...\n"))
		ctx.Writer.Flush()

		classify_response, err := first_parser.SendRequest(username, string(req.age), req.join_reason, req.about, req.inviter, context.Background(), gemini)
		if err != nil {
			ctx.JSON(422, "Говно ответ от gemini, "+err.Error())
			return
		}

		classify_data := first_parser.ParseSomeDataData(classify_response.Text())

		ctx.Writer.Write([]byte("Выносим финальный вердикт...\n"))
		ctx.Writer.Flush()

		verdict_response, err := final_reviewer.SendRequest(classify_data, username, string(req.age), req.join_reason, req.about, req.inviter, context.Background(), gemini)
		if err != nil {
			ctx.JSON(422, "Говно ответ от gemini, "+err.Error())
			return
		}

		var verdict_json map[string]string
		err = json.Unmarshal([]byte(verdict_response.Text()), &verdict_json)
		if err != nil {
			ctx.JSON(422, "Говно ответ от gemini, "+err.Error())
			return
		}

		fmt.Printf("\n\n%v\n", verdict_json)
		ctx.Writer.Write([]byte(fmt.Sprintf("{\"action\": \"%s\", \"answer\": \"%s\"}\n", verdict_json["action"], verdict_json["reason"])))
	})

	log.Printf("Запуск Gin спустя: %s с начала запуска программы", time.Since(appStart))
	r.Run(":" + listenPort)
}
