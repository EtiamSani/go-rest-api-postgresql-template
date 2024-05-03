package main

import (
	"fmt"
	"os"

	"github.com/etiamsani/go-rest-api-postgresl-template/api/route"
	"github.com/etiamsani/go-rest-api-postgresl-template/api/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {
	store.LoadEnvVariables()
	store.ConnectToDB()
}


func main() {
	err:= godotenv.Load()
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	router := route.UseRouter()

    router.Run("localhost:" + os.Getenv("PORT"))
	logger,_ := zap.NewProduction()
  	defer logger.Sync()
	logger.Sugar().Info("🟢 Server is running on port", os.Getenv("PORT"))
}