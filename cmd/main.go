package main

import (
	"log"

	"eth_tracker/internal/app"
	"eth_tracker/internal/config"
	"eth_tracker/pkg/db"
	"eth_tracker/pkg/eth"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := config.Load()
	dbConn, err := db.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	ethClient := eth.NewEthClient(cfg.APIKey)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}

	application := app.NewApplication(dbConn, ethClient)
	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
