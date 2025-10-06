package utils

import (
	"log"

	"github.com/joho/godotenv"
)

// loadEnv Load env variable(s) from .env file
func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Fail loading .env: %s", err)
	}
}
