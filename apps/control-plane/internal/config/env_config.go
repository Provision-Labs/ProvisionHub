package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port           int      `env:"PORT" default:"8180"`
	Env            string   `env:"ENV" default:"development"`
	Issuer         string   `env:"ISSUER" default:""`
	ClientId       string   `env:"CLIENT_ID" default:""`
	ClientSecret   string   `env:"CLIENT_SECRET" default:""`
	RedirectURL    string   `env:"REDIRECT_URL" default:""`
	Scopes         []string `env:"SCOPES" default:"" split:","`
	LogoutRedirect string   `env:"LOGOUT_REDIRECT_URL" default:""`
	SessionSecret  string   `env:"SESSION_SECRET" default:""`
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
