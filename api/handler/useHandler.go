package handler

import (
	"net/http"

	"github.com/etiamsani/go-rest-api-postgresl-template/api/model"
	"github.com/etiamsani/go-rest-api-postgresl-template/api/store"
	"github.com/gin-gonic/gin"
)

func UseHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Hello World"})
}

func GetAllUsers(c *gin.Context) {
	
}

func CreateUser(c *gin.Context) {
	user := model.User{Name: "Jinzhu", Email: "Jinzhu@example.com"}

	result := store.DB.Create(&user) 

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{
		"user": user,
	})
}