package main

import (
	"github.com/etiamsani/go-rest-api-postgresl-template/api/model"
	"github.com/etiamsani/go-rest-api-postgresl-template/api/store"
)

func init() {
	store.LoadEnvVariables()
	store.ConnectToDB()
}

func main() {

	store.DB.AutoMigrate(&model.User{})
}