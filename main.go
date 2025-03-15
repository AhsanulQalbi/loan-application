package main

import (
	"log"
	"os"

	"loan_application/config"
	"loan_application/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := config.ConnectDB()
	routes := routes.SetupRouter(db)
	log.Fatal(routes.Run(":" + os.Getenv("APP_PORT")))
}
