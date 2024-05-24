package route

import (
	"github.com/etiamsani/go-rest-api-postgresl-template/api/handler"
	"github.com/etiamsani/go-rest-api-postgresl-template/api/pkg/middleware"
	"github.com/gin-gonic/gin"
)


func OauthRouter(r *gin.Engine) {
	
	r.GET("/auth/:provider/callback", handler.GetAuthCallBackFunction)
	
	r.GET("/logout/:provider", handler.LogoutHandler)
	
	r.GET("/auth/:provider", handler.AuthHandler)



}

func StripeRouter( r *gin.Engine)  {
	r.POST("/stripe/checkout", handler.CheckoutCreator) 
	// r.GET("/stripe/webhook" , handler.HandleEvent)
}

func UserRouter( r *gin.Engine)  {
	r.GET("/use", handler.GetAllUsers) 
	r.GET("/user/me", middleware.JWTAuthMiddleware(), middleware.OAuthAuthMiddleware(), handler.GetUserData)
	r.POST("/user/signup", handler.Signup)
	r.POST("/user/login", handler.Login)
}