package main

import (
	"log"

	"github.com/andreyxaxa/auth-svc/config"
	"github.com/andreyxaxa/auth-svc/internal/app"
	"github.com/joho/godotenv"
)

// @title   Auth service API
// @version 1.0
// @host    localhost:8080
// BasePath /v1

// @securityDefinitions.apikey TokenAuth
// @in                         header
// @name 					   Authorization

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	// Config
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error :%s", err)
	}

	// Run
	app.Run(cfg)
}
