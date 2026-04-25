package config

import (
	"os"
	"strconv"
	"time"

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
		Env     string
	}

	HTTP struct {
		Port string
	}

	Log struct {
		Level string
	}

	PG struct {
		URL             string
		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime time.Duration
	}
)

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		App: App{
			Name:    os.Getenv("APP_NAME"),
			Version: os.Getenv("APP_VERSION"),
			Env:     os.Getenv("ENV"),
		},
		HTTP: HTTP{
			Port: os.Getenv("HTTP_PORT"),
		},
		Log: Log{
			Level: os.Getenv("LOG_LEVEL"),
		},
		PG: PG{
			URL:             os.Getenv("PG_URL"),
			MaxOpenConns:    getEnvInt("PG_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    getEnvInt("PG_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: time.Duration(getEnvInt("PG_CONN_MAX_LIFETIME_SEC", 3600)) * time.Second,
		},
	}

	return cfg, nil
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultVal
}
