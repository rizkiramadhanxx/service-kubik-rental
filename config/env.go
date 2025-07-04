package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	JWTSecret []byte
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is missing in .env file")
	}

	JWTSecret = []byte(secret)
}
