package api

import (
	"log"
	"os"

	"github.com/norfabagas/auth/api/controllers"
)

var server = controllers.Server{}

func Run() {
	var err error

	if err != nil {
		log.Fatalf("Error getting env: %v", err)
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	if apPKey := os.Getenv("APP_KEY"); len([]rune(apPKey)) != 32 {
		log.Fatalf("APP_KEY is not 32 bit long")
	}

	server.Run(":" + os.Getenv("PORT"))
}
