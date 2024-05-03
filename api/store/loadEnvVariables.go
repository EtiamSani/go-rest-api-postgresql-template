package store

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		fmt.Print("Error loading .env file")
	}
}