package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

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
    ID        uint       `gorm:"primarykey" json:"id"`
    CreatedAt time.Time  `json:"-"`
    UpdatedAt time.Time  `json:"-"`
    Contents  string `gorm:"not null" json:"contents"`
}

func main() {
	engine := gin.Default()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("fail to connect database.")
	}

	db.AutoMigrate(&Todo{})

	engine.GET("/todo", func(c *gin.Context) {
		var todos []Todo

		db.Find(&todos)

		c.JSON(http.StatusOK, gin.H{
			"items": todos,
		})
	})

	engine.Run(":8000")
}
