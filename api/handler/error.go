package handler

import (
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, statusCode int, errorMessage string) {
    c.JSON(statusCode, gin.H{"error": errorMessage})
}
