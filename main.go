package main

import (
	"fmt"
	"os"

	"github.com/etiamsani/go-rest-api-postgresl-template/api/route"
	"github.com/etiamsani/go-rest-api-postgresl-template/api/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v78"
	"go.uber.org/zap"
)

func init() {
	store.LoadEnvVariables()
	store.ConnectToDB()
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}


func main() {
	
	err:= godotenv.Load()
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	corsConfig := cors.DefaultConfig()
    corsConfig.AllowOrigins = []string{"http://localhost:5173"} 
    corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"} 
    corsConfig.AllowHeaders = []string{"Content-Type", "Authorization", "X-Requested-With", "Accept", "Origin", "Access-Control-Allow-Origin"}
	corsConfig.AllowCredentials = true

	router := gin.Default()
	store.NewAuth()

	router.Use(cors.New(corsConfig))

	route.OauthRouter(router)
	route.UserRouter(router)
	route.StripeRouter(router)

    router.Run("localhost:" + os.Getenv("PORT"))
	logger,_ := zap.NewProduction()
  	defer logger.Sync()
	logger.Sugar().Info("🟢 Server is running on port", os.Getenv("PORT"))
}