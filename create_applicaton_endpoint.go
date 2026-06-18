package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"svc-applications/final_reviewer"
	"svc-applications/first_parser"

	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
)

func createApplicationHandler(ctx *gin.Context) {
	//db := ctx.MustGet("db").(*sql.DB)
	gemini := ctx.MustGet("gemini").(*genai.Client)
	usersServiceUrl := ctx.MustGet("usersServiceUrl").(string)
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

	classify_response, err := first_parser.SendRequest(username, string(req.age), req.about, req.inviter, ctx.Request.Context(), gemini)
	if err != nil {
		ctx.JSON(422, "Говно ответ от gemini, "+err.Error())
		return
	}

	classify_data := first_parser.ParseSomeDataData(classify_response.Text())

	ctx.Writer.Write([]byte("Выносим финальный вердикт...\n"))
	ctx.Writer.Flush()

	verdict_response, err := final_reviewer.SendRequest(classify_data, username, string(req.age), req.about, req.inviter, ctx.Request.Context(), gemini)
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

}
