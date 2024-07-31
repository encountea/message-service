package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Kafka    KafkaConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	config := &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("SSL_MODE"),
		},
		Kafka: KafkaConfig{
			Brokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
			Topic:   os.Getenv("KAFKA_TOPIC"),
		},
	}

	return config, nil
}
