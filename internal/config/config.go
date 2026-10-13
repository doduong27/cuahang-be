package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	Environment string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=192.168.1.253 user=app_cuahang_dev password=DevPass@2026 dbname=cuahang_dev port=5432 sslmode=disable"
	}

	env := os.Getenv("NODE_ENV")
	if env == "" {
		env = os.Getenv("ENV")
		if env == "" {
			env = "dev"
		}
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		Environment: env,
	}
}
