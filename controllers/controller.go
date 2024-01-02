package controllers

import (
	"gin-sample/api"
	"gin-sample/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TodoController struct {
	db *gorm.DB
}

func NewTodoController(db *gorm.DB) *TodoController {
	return &TodoController{db: db}
}

func (c *TodoController) GetTodoListHandler(context *gin.Context) {
	var todos []models.Todo

	c.db.Order("id").Find(&todos)

	context.JSON(http.StatusOK, gin.H{
		"items": todos,
	})
}

func (c *TodoController) GetTodoHandler(context *gin.Context) {
	todo := models.Todo{}

	id := context.Param("id")
	c.db.First(&todo, id)

	context.JSON(http.StatusOK, gin.H{
		"item": todo,
	})
}

func (c *TodoController) PostTodoHandler(context *gin.Context) {
	var data api.EditTodoRequest

	if err := context.BindJSON(&data); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid params",
		})
		return
	}

	todo := models.Todo{}
	c.db.First(&todo, data.ID)

	todo.Contents = data.Contents
	c.db.Save(&todo)

	context.Status(http.StatusOK)
}

func (c *TodoController) PutTodoHandler(context *gin.Context) {
	var data api.CreateTodoRequest

	if err := context.BindJSON(&data); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid params",
		})
		return
	}

	c.db.Create(&models.Todo{Contents: data.Contents})

	context.Status(http.StatusCreated)
}

func (c *TodoController) DeleteTodoHandler(context *gin.Context) {
	var data api.DeleteTodoRequest

	if err := context.BindJSON(&data); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid params",
		})
		return
	}

	todo := models.Todo{}
	c.db.First(&todo, data.ID)
	c.db.Delete(&todo)

	context.Status(http.StatusNoContent)
}
