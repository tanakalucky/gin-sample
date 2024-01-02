package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gin-sample/controllers"
	"gin-sample/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host        = os.Getenv("HOST")
	username    = os.Getenv("USERNAME")
	password    = os.Getenv("USERPASS")
	database    = os.Getenv("DATABASE")
	dsn         = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Tokyo", host, username, password, database)
	allowOrigin = os.Getenv("ALLOWORIGIN")
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			allowOrigin,
		},
		AllowMethods: []string{
			"POST",
			"GET",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		MaxAge: 24 * time.Hour,
	}))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("fail to connect database.")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get sql.DB from Gorm")
	}
	defer sqlDB.Close()

	db.AutoMigrate(&models.Todo{})

	con := controllers.NewTodoController(db)

	r.GET("/todo", con.GetTodoListHandler)
	r.GET("/todo/:id", con.GetTodoHandler)
	r.POST("/todo/edit", con.PostTodoHandler)
	r.PUT("/todo", con.PutTodoHandler)
	r.DELETE("/todo/delete", con.DeleteTodoHandler)
	// プリフライトリクエスト用
	r.OPTIONS("/todo", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	r.Run(":8000")
}
