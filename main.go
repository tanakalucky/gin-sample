package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gin-sample/api"
	"gin-sample/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = os.Getenv("HOST")
	username = os.Getenv("USERNAME")
	password = os.Getenv("USERPASS")
	database = os.Getenv("DATABASE")
	dsn      = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Tokyo", host, username, password, database)
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://host.docker.internal",
			"http://localhost",
			"http://0.0.0.0",
			"*",
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

	db.AutoMigrate(&models.Todo{})

	r.GET("/todo", func(c *gin.Context) {
		var todos []models.Todo

		db.Order("id").Find(&todos)

		c.JSON(http.StatusOK, gin.H{
			"items": todos,
		})
	})

	r.OPTIONS("/todo", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	r.PUT("/todo", func(c *gin.Context) {
		var data api.CreateTodoRequest

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid params",
			})
			return
		}

		db.Create(&models.Todo{Contents: data.Contents})

		c.Status(http.StatusCreated)
	})

	r.DELETE("/todo/delete", func(c *gin.Context) {
		var data api.DeleteTodoRequest

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid params",
			})
			return
		}

		todo := models.Todo{}
		db.First(&todo, data.ID)
		db.Delete(&todo)

		c.Status(http.StatusNoContent)
	})

	r.GET("/todo/:id", func(c *gin.Context) {
		todo := models.Todo{}

		id := c.Param("id")
		db.First(&todo, id)

		c.JSON(http.StatusOK, gin.H{
			"item": todo,
		})
	})

	r.POST("/todo/edit", func(c *gin.Context) {
		var data api.EditTodoRequest

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid params",
			})
			return
		}

		todo := models.Todo{}
		db.First(&todo, data.ID)

		todo.Contents = data.Contents
		db.Save(&todo)

		c.Status(http.StatusOK)
	})

	r.Run(":8000")
}
