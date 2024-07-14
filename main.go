package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	route.VerifyRouter(router)



    // router.Run("localhost:" + os.Getenv("PORT"))
	server := &http.Server{
		Addr:    "localhost:" + os.Getenv("PORT"),
		Handler: router,
	}

	
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Ajouter un WaitGroup pour attendre que les goroutines se terminent
	var wg sync.WaitGroup
	wg.Add(2)

	

	// Démarrer le serveur HTTP dans une goroutine
	go func() {
		defer wg.Done()
		logger,_ := zap.NewProduction()
		defer logger.Sync()
		logger.Sugar().Info("🟢 Server is running on port", os.Getenv("PORT"))

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		// Attendre un signal d'interruption ou SIGTERM
		<-ctx.Done()

		// Fermer le serveur gracieusement
		log.Println("Shutting down server...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		  }
		
	}()

	

	// // Créer un contexte pour la fermeture du serveur
	// shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// // Arrêter le serveur
	// if err := server.Shutdown(shutdownCtx); err != nil {
	// 	log.Fatalf("Server shutdown failed: %v", err)
	// }

	// Attendre que toutes les goroutines se terminent
	wg.Wait()

	log.Println("Server has been gracefully shutdown")
}