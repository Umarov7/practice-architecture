package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
	"go.uber.org/fx"
)

var Module = fx.Options(fx.Provide(Load))

type Config struct {
	ADDRESS      string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration

	// Postgres
	Postgres_HOST     string
	Postgres_PORT     string
	Postgres_USER     string
	Postgres_NAME     string
	Postgres_PASSWORD string

	// MongoDB
	MongoDB_URI        string
	MongoDB_NAME       string
	MongoDB_COLLECTION string
}

func Load() *Config {
	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("error while getting current working directory: %v", err)
	}

	if err := godotenv.Load(path + "/.env"); err != nil {
		log.Printf("error while loading .env file: %v", err)
	}

	return &Config{
		ADDRESS: cast.ToString(coalesce("ADDRESS", "localhost:8080")),

		WriteTimeout: cast.ToDuration(coalesce("WRITE_TIMEOUT", "5s")),
		ReadTimeout:  cast.ToDuration(coalesce("READ_TIMEOUT", "5s")),

		// Postgres
		Postgres_HOST:     cast.ToString(coalesce("POSTGRES_HOST", "localhost")),
		Postgres_PORT:     cast.ToString(coalesce("POSTGRES_PORT", "5432")),
		Postgres_USER:     cast.ToString(coalesce("POSTGRES_USER", "postgres")),
		Postgres_NAME:     cast.ToString(coalesce("POSTGRES_NAME", "postgres")),
		Postgres_PASSWORD: cast.ToString(coalesce("POSTGRES_PASSWORD", "")),

		// MongoDB
		MongoDB_URI:        cast.ToString(coalesce("MONGO_DB_URI", "")),
		MongoDB_NAME:       cast.ToString(coalesce("MONGO_DB_NAME", "")),
		MongoDB_COLLECTION: cast.ToString(coalesce("MONGO_DB_COLLECTION", "")),
	}
}

func coalesce(key string, value interface{}) interface{} {
	val, exists := os.LookupEnv(key)
	if exists {
		return val
	}
	return value
}
