package main

import (
	"log"

	"github.com/etiamsani/go-rest-api-postgresl-template/api/model"
	"github.com/etiamsani/go-rest-api-postgresl-template/api/store"
)

func init() {
	store.LoadEnvVariables()
	store.ConnectToDB()
}

func main() {

	log.Println("Démarrage des migrations...")

	// Migration pour le modèle User
	if err := store.DB.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Échec de la migration pour User: %v", err)
	} else {
		log.Println("Migration réussie pour User.")
	}

	// Migration pour le modèle VerificationData
	if err := store.DB.AutoMigrate(&model.VerificationData{}); err != nil {
		log.Fatalf("Échec de la migration pour VerificationData: %v", err)
	} else {
		log.Println("Migration réussie pour VerificationData.")
	}

	log.Println("Toutes les migrations se sont bien déroulées.")
}