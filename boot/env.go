package boot

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvFile() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
		return
	}

	log.Println("Loading .env file successfully")
}
