package config

import (
	httprouter "api-for-people/internal/transport/http-router"
	"api-for-people/pkg/postgres"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

type Config struct {
	PostgresConfig postgres.Config
	RouterConfig   httprouter.Config
}

func NewConfig() *Config {
	slog.Info("Loading config")
	defer slog.Info("Config Loaded.")
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Error loading .env file")
	}

	return &Config{
		PostgresConfig: postgres.Config{
			Host:     os.Getenv("HOST"),
			Port:     os.Getenv("PORT"),
			User:     os.Getenv("USER"),
			Password: os.Getenv("PASSWORD"),
			Database: os.Getenv("DATABASE"),
		},
		RouterConfig: httprouter.Config{
			Address: os.Getenv("ADDRESS"),
		},
	}
}
