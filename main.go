package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	config "github.com/punchanabu/redrice-backend-go/config"
	routers "github.com/punchanabu/redrice-backend-go/routers"
	"github.com/punchanabu/redrice-backend-go/routers/api"
	v1 "github.com/punchanabu/redrice-backend-go/routers/api/v1"
)

func main() {

	// Load the env
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Setup database connection
	db := config.SetupDBConnection()
	if db == nil {
		log.Fatal("Failed to connect to database!")
	}

	// Initialize necessary handlers
	v1.InitializedUserHandler(db)
	v1.InitializedRestaurantHandler(db)
	api.InitializedAuthHandler(db)

	// Initialize router
	r := routers.UseRouter()
	r.Use(cors.New(config.CORSConfig()))

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal("Server failed to start!")
	}

}
