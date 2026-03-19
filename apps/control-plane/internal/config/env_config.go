package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port int    `envconfig:"PORT" default:"8180"`
	Env  string `envconfig:"ENV" default:"development"`

	// Keycloak vars
	Issuer         string   `envconfig:"ISSUER" required:"true"`
	ClientId       string   `envconfig:"CLIENT_ID" required:"true"`
	ClientSecret   string   `envconfig:"CLIENT_SECRET" required:"true"`
	RedirectURL    string   `envconfig:"REDIRECT_URL" required:"true"`
	Scopes         []string `envconfig:"SCOPES" required:"true"`
	LogoutRedirect string   `envconfig:"LOGOUT_REDIRECT_URL" required:"true"`
	SessionSecret  string   `envconfig:"SESSION_SECRET" required:"true"`

	// Database vars
	Host      string `envconfig:"DB_HOST" required:"true"`
	PortDB    int    `envconfig:"DB_PORT" default:"5432"`
	Username  string `envconfig:"DB_USERNAME" required:"true"`
	Password  string `envconfig:"DB_PASSWORD" required:"true"`
	DbName    string `envconfig:"DB_NAME" required:"true"`
	DbSSLMode string `envconfig:"DB_SSL_MODE" default:"disable"`
}

var cfg Config

// LoadConfig loads environment variables into cfg
func LoadConfig() *Config {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	if env != "production" {
		filename := ".env." + env // → ".env.development"
		if err := godotenv.Load(filename); err != nil {
			log.Printf("No %s file found, relying on environment variables\n", filename)
		}
	}

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to load env vars: %v", err)
	}
	return &cfg
}
