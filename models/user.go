package models

import (
	"gorm.io/gorm"
)

type User struct {
	UserID    uint   `gorm:"primaryKey" json:"userId"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
	Role      string `json:"role"`
	Password  string `json:"password"`
	gorm.Model
}
