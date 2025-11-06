package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvFile() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Loading .env file successfully")
}
