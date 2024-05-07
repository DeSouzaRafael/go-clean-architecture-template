package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		App
		HTTP
		Log
		PG
	}

	App struct {
		Name    string
		Version string
	}

	HTTP struct {
		Port string
	}

	Log struct {
		Level string
	}

	PG struct {
		URL string
	}
)

func NewConfig() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := &Config{
		App: App{
			Name:    os.Getenv("APP_NAME"),
			Version: os.Getenv("APP_VERSION"),
		},
		HTTP: HTTP{
			Port: os.Getenv("HTTP_PORT"),
		},
		Log: Log{
			Level: os.Getenv("LOG_LEVEL"),
		},
		PG: PG{
			URL: os.Getenv("PG_URL"),
		},
	}

	return cfg, nil
}
