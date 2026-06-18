package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

func getAllHandler(ctx *gin.Context) {
	db := ctx.MustGet("db").(*sql.DB)

	rows, err := db.QueryContext(ctx.Request.Context(), "SELECT id,user_id,user_ip,age,about,inviter,ai_categories,ai_decision,ai_answer,ai_comment,admin_decision,admin_id FROM applications")
	if err != nil {
		log.Panicln(err)
	}

	var data []ApplicationData = make([]ApplicationData, 0)
	for rows.Next() {
		var dat ApplicationData
		if err := rows.Scan(&dat.id, &dat.user_id, &dat.user_ip, &dat.age, &dat.about, &dat.inviter, &dat.ai_categories, &dat.ai_decision, &dat.ai_answer, &dat.ai_comment, &dat.admin_decision, &dat.admin_id); err != nil {
			log.Panicln(err)
		}
		data = append(data, dat)
	}

	if err := rows.Err(); err != nil {
		log.Panicln(err)
	}

	ctx.JSON(200, data)
}

func getByUUIDHandler(ctx *gin.Context) {
	db := ctx.MustGet("db").(*sql.DB)
	uuid := ctx.Param("uuid")

	rows, err := db.QueryContext(ctx.Request.Context(), "SELECT id,user_id,user_ip,age,about,inviter,ai_categories,ai_decision,ai_answer,ai_comment,admin_decision,admin_id FROM applications WHERE user_id = ?", uuid)
	if err != nil {
		log.Panicln(err)
	}

	var data []ApplicationData = make([]ApplicationData, 0)
	for rows.Next() {
		var dat ApplicationData
		if err := rows.Scan(&dat.id, &dat.user_id, &dat.user_ip, &dat.age, &dat.about, &dat.inviter, &dat.ai_categories, &dat.ai_decision, &dat.ai_answer, &dat.ai_comment, &dat.admin_decision, &dat.admin_id); err != nil {
			log.Panicln(err)
		}
		data = append(data, dat)
	}

	if err := rows.Err(); err != nil {
		log.Panicln(err)
	}

	ctx.JSON(200, data)
}

func getMyHandler(ctx *gin.Context) {
	db := ctx.MustGet("db").(*sql.DB)
	user_id := ctx.Request.Header.Get("eauth-user-id")

	rows, err := db.QueryContext(ctx.Request.Context(), "SELECT id,user_id,age,about,inviter,ai_categories,ai_decision,ai_answer,admin_decision,admin_id FROM applications WHERE user_id = ?", user_id)
	if err != nil {
		log.Panicln(err)
	}

	var data []ApplicationDataShort = make([]ApplicationDataShort, 0)
	for rows.Next() {
		var dat ApplicationDataShort
		if err := rows.Scan(&dat.id, &dat.user_id, &dat.age, &dat.about, &dat.inviter, &dat.ai_categories, &dat.ai_decision, &dat.ai_answer, &dat.admin_decision, &dat.admin_id); err != nil {
			log.Panicln(err)
		}
		data = append(data, dat)
	}

	if err := rows.Err(); err != nil {
		log.Panicln(err)
	}

	ctx.JSON(200, data)
}
