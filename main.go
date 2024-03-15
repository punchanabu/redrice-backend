package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/punchanabu/redrice-backend-go/models"
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

}
