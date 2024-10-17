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

	// Kafka
	KAFKA_ADDRESS                string
	KAFKA_TOPIC_USER_CREATED     string
	KAFKA_TOPIC_USER_UPDATED     string
	KAFKA_TOPIC_USER_DELETED     string
	KAFKA_TOPIC_COMPUTER_CREATED string
	KAFKA_TOPIC_COMPUTER_UPDATED string
	KAFKA_TOPIC_COMPUTER_DELETED string

	// RabbitMQ
	RabbitMQ_ADDRESS                string
	RabbitMQ_QUEUE_USER_CREATED     string
	RabbitMQ_QUEUE_USER_UPDATED     string
	RabbitMQ_QUEUE_USER_DELETED     string
	RabbitMQ_QUEUE_COMPUTER_CREATED string
	RabbitMQ_QUEUE_COMPUTER_UPDATED string
	RabbitMQ_QUEUE_COMPUTER_DELETED string
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

		// Kafka
		KAFKA_ADDRESS:                cast.ToString(coalesce("KAFKA_ADDRESS", "localhost:9092")),
		KAFKA_TOPIC_USER_CREATED:     cast.ToString(coalesce("KAFKA_TOPIC_USER_CREATED", "USER_CREATED")),
		KAFKA_TOPIC_USER_UPDATED:     cast.ToString(coalesce("KAFKA_TOPIC_USER_UPDATED", "USER_UPDATED")),
		KAFKA_TOPIC_USER_DELETED:     cast.ToString(coalesce("KAFKA_TOPIC_USER_DELETED", "USER_DELETED")),
		KAFKA_TOPIC_COMPUTER_CREATED: cast.ToString(coalesce("KAFKA_TOPIC_COMPUTER_CREATED", "COMPUTER_CREATED")),
		KAFKA_TOPIC_COMPUTER_UPDATED: cast.ToString(coalesce("KAFKA_TOPIC_COMPUTER_UPDATED", "COMPUTER_UPDATED")),
		KAFKA_TOPIC_COMPUTER_DELETED: cast.ToString(coalesce("KAFKA_TOPIC_COMPUTER_DELETED", "COMPUTER_DELETED")),

		// RabbitMQ
		RabbitMQ_ADDRESS:                cast.ToString(coalesce("RabbitMQ_ADDRESS", "amqp://guest:guest@localhost:5672")),
		RabbitMQ_QUEUE_USER_CREATED:     cast.ToString(coalesce("RabbitMQ_QUEUE_USER_CREATED", "USER_CREATED")),
		RabbitMQ_QUEUE_USER_UPDATED:     cast.ToString(coalesce("RabbitMQ_QUEUE_USER_UPDATED", "USER_UPDATED")),
		RabbitMQ_QUEUE_USER_DELETED:     cast.ToString(coalesce("RabbitMQ_QUEUE_USER_DELETED", "USER_DELETED")),
		RabbitMQ_QUEUE_COMPUTER_CREATED: cast.ToString(coalesce("RabbitMQ_QUEUE_COMPUTER_CREATED", "COMPUTER_CREATED")),
		RabbitMQ_QUEUE_COMPUTER_UPDATED: cast.ToString(coalesce("RabbitMQ_QUEUE_COMPUTER_UPDATED", "COMPUTER_UPDATED")),
		RabbitMQ_QUEUE_COMPUTER_DELETED: cast.ToString(coalesce("RabbitMQ_QUEUE_COMPUTER_DELETED", "COMPUTER_DELETED")),
	}
}

func coalesce(key string, value interface{}) interface{} {
	val, exists := os.LookupEnv(key)
	if exists {
		return val
	}
	return value
}
