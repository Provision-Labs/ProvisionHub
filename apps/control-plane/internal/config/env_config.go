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
	PluginsRegistryPath string `envconfig:"PLUGINS_REGISTRY_PATH" default:"plugins.json"`

	// Auth plugin connection vars
	AuthPluginAddr     string `envconfig:"AUTH_PLUGIN_ADDR" default:"127.0.0.1:50051"`
	AuthPluginInsecure bool   `envconfig:"AUTH_PLUGIN_INSECURE" default:"true"`
	AuthPluginTimeoutMs int   `envconfig:"AUTH_PLUGIN_TIMEOUT_MS" default:"2000"`

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
