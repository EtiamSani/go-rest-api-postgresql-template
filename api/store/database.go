package store

import (
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

  var DB *gorm.DB

  

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	logger,_ := zap.NewProduction()
  	defer logger.Sync()
	logger.Sugar().Info("🟢 Connection Opened to Database")
}