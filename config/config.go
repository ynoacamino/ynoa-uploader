package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var CLOUDINARY_URL string = ""
var DATABASE_URL string = ""
var SECRET_CONNECTION string = ""

func InitEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	CLOUDINARY_URL = os.Getenv("CLOUDINARY_URL")
	if CLOUDINARY_URL == "" {
		log.Fatal("CLOUDINARY_URL is not set")
	}

	DATABASE_URL = os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	SECRET_CONNECTION = os.Getenv("SECRET_CONNECTION")
	if SECRET_CONNECTION == "" {
		log.Fatal("SECRET_CONNECTION is not set")
	}
}
