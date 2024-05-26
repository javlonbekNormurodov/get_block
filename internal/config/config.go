package config

import (
	"os"
)

type Config struct {
	APIKey      string
	DatabaseURL string
}

func Load() Config {
	return Config{
		APIKey:      os.Getenv("API_KEY"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}
