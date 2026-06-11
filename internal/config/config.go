package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string `env:"JWT_SECRET"`
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env"); err == nil {
		log.Println("Loaded configuration from .env")
	} else if err := godotenv.Load(".env.docker"); err == nil {
		log.Println("Loaded configuration from .env.docker")
	} else {
		log.Println("No env file found, relying on environment variables")
	}

	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),

		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
