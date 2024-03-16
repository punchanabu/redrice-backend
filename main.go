package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/punchanabu/redrice-backend-go/models"
	"github.com/punchanabu/redrice-backend-go/routers"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := gorm.Open(mysql.Open("DB_CONN"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	// auto migration
	db.AutoMigrate(&models.User{}, &models.Restaurant{}, &models.Reservation{})

	r := routers.UseRouter()
	if err := r.Run(os.Getenv("PORT")); err != nil {
		log.Fatal("Server failed to start!")
	}

}
