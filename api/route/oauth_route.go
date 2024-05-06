package route

import (
	"github.com/etiamsani/go-rest-api-postgresl-template/api/handler"
	"github.com/gin-gonic/gin"
)

func OauthRouter(r *gin.Engine) {
	
	r.GET("/auth/:provider/callback", handler.GetAuthCallBackFunction)
	
	r.GET("/logout/:provider", handler.LogoutHandler)
	
	r.GET("/auth/:provider", handler.AuthHandler)



}