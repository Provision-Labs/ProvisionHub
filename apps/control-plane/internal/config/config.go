package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Env  string
}

func Load() *Config {
	// Only loads .env in development environment
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Printf("Error loading .env file: %v", err)
		}
	}

	config := &Config{
		Port: getEnv("PORT", "8180"),
		Env:  getEnv("ENV", "development"),
	}

	log.Printf("Config is loaded from env: %v", config.Env)

	return config
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
