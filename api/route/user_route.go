package route

import (
	"github.com/etiamsani/go-rest-api-postgresl-template/api/handler"
	"github.com/etiamsani/go-rest-api-postgresl-template/api/middleware"
	"github.com/gin-gonic/gin"
)

func UseRouter( r *gin.Engine)  {
	r.GET("/use", handler.GetAllUsers) 
	r.GET("/user/me", middleware.AuthMiddleware(), handler.GetUserData)
}