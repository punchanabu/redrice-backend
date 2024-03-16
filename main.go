package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/punchanabu/redrice-backend-go/routers"
	"github.com/gin-contrib/cors"
	"github.com/punchanabu/redrice-backend-go/config"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := config.SetupDBConnection()
	if db == nil {
		log.Fatal("Failed to connect to database!")
	}

	r := routers.UseRouter()

	r.Use(cors.New(config.CORSConfig()))

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal("Server failed to start!")
	}

}
