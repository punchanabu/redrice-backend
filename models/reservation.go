package models

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	DateTime     time.Time `json:"dateTime"`
	UserID       uint      `json:"userId"`
	User         User      `gorm:"foreignKey:UserID" json:"user"`
	RestaurantID uint      `json:"restaurantId"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant"`
	gorm.Model
}
