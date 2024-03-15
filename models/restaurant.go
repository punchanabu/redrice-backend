package models

import (
	"gorm.io/gorm"
)

type Restaurant struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	Description string `json:"description"`
	gorm.Model
}
