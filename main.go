package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

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

type Todo struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Contents  string    `gorm:"not null" json:"contents"`
}

type CreateTodoRequest struct {
	Contents string `json:"contents"`
}

type DeleteTodoRequest struct {
	ID int `json:"id"`
}

type EditTodoRequest struct {
	ID       int    `json:"id"`
	Contents string `json:"contents"`
}

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

	db.AutoMigrate(&Todo{})

	r.GET("/todo", func(c *gin.Context) {
		var todos []Todo

		db.Order("id").Find(&todos)

		c.JSON(http.StatusOK, gin.H{
			"items": todos,
		})
	})

	r.OPTIONS("/todo", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	r.PUT("/todo", func(c *gin.Context) {
		var data CreateTodoRequest

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid params",
			})
			return
		}

		db.Create(&Todo{Contents: data.Contents})

		c.Status(http.StatusCreated)
	})

	r.DELETE("/todo/delete", func(c *gin.Context) {
		var data DeleteTodoRequest

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid params",
			})
			return
		}

		todo := Todo{}
		db.First(&todo, data.ID)
		db.Delete(&todo)

		c.Status(http.StatusNoContent)
	})

	r.GET("/todo/:id", func(c *gin.Context) {
		todo := Todo{}

		id := c.Param("id")
		db.First(&todo, id)

		c.JSON(http.StatusOK, gin.H{
			"item": todo,
		})
	})

	r.POST("/todo/edit", func(c *gin.Context) {
		var data EditTodoRequest

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid params",
			})
			return
		}

		todo := Todo{}
		db.First(&todo, data.ID)

		todo.Contents = data.Contents
		db.Save(&todo)

		c.Status(http.StatusOK)
	})

	r.Run(":8000")
}
