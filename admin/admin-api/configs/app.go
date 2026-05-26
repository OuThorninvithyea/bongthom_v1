package config

import (
	// Community Packages
	"log"
	"os"
	"github.com/joho/godotenv"

	// Internal Packages
	"admin-api/pkg/utls"
)

type AppConfig struct {
	AppHost string
	AppPort int
}

func NewConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error Loading .env file %v", err)
	}
	host := os.Getenv("API_HOST")
	port := utls.GetenvInt("API_PORT", 8888)
	return &AppConfig{
		AppHost: host,
		AppPort: port,
	}
}
