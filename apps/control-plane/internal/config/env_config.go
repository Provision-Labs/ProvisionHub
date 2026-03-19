package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port int    `env:"PORT" default:"8180"`
	Env  string `env:"ENV" default:"development"`

	// Keycloak vars
	Issuer         string   `env:"ISSUER" required:"true"`
	ClientId       string   `env:"CLIENT_ID" required:"true"`
	ClientSecret   string   `env:"CLIENT_SECRET" required:"true"`
	RedirectURL    string   `env:"REDIRECT_URL" required:"true"`
	Scopes         []string `env:"SCOPES" required:"true" split:","`
	LogoutRedirect string   `env:"LOGOUT_REDIRECT_URL" required:"true"`
	SessionSecret  string   `env:"SESSION_SECRET" required:"true"`

	// Database vars
	Host     string `env:"DB_HOST" required:"true"`
	Username string `env:"DB_USERNAME" required:"true"`
	Password string `env:"DB_PASSWORD" required:"true"`
	DbName   string `env:"DB_NAME" required:"true"`
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
