package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	config "github.com/punchanabu/redrice-backend-go/config"
	routers "github.com/punchanabu/redrice-backend-go/routers"
	v1 "github.com/punchanabu/redrice-backend-go/routers/api/v1"
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

	v1.InitializedUserHandler(db)

	r := routers.UseRouter()
	r.Use(cors.New(config.CORSConfig()))

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal("Server failed to start!")
	}

}
