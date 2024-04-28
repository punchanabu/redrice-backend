package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"unique"`
	Telephone    string `json:"telephone" gorm:"unique"`
	Role         string `json:"role"`
	Password     string `json:"password"`
	RestaurantId uint   `json:"restaurant_id"`
	gorm.Model   `json:"-" swaggerignore:"true"`
}

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db}
}

func (h UserHandler) CreateUser(user *User) error {
	// Check if email already exists
	existingEmail, err := h.GetUserByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingEmail != nil {
		return fmt.Errorf("email already exists")
	}

	// Check if telephone already exists
	existingTelephone, err := h.GetUserByTelephone(user.Telephone)
	if err != nil {
		return err
	}
	if existingTelephone != nil {
		return fmt.Errorf("telephone already exists")
	}

	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Create the user
	return h.db.Create(user).Error
}

func (h *UserHandler) CheckPassword(email, password string) bool {
	var user User
	if err := h.db.Where("email = ?", email).First(&user).Error; err != nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (h *UserHandler) GetUser(id uint) (*User, error) {
	var user User
	result := h.db.First(&user, id)
	return &user, result.Error
}

func (h *UserHandler) GetUsers() ([]User, error) {
	var users []User
	result := h.db.Find(&users)
	return users, result.Error
}

func (h *UserHandler) UpdateUser(id uint, user *User) error {
	result := h.db.Model(&User{}).Where("id = ?", id).Updates(user)
	return result.Error
}

func (h *UserHandler) DeleteUser(id uint) error {
	result := h.db.Delete(&User{}, id)
	return result.Error
}

func (h *UserHandler) GetUserByEmail(email string) (*User, error) {
	var user User
	result := h.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, result.Error
}

func (h *UserHandler) GetUserByTelephone(telephone string) (*User, error) {
	var user User
	result := h.db.Where("telephone = ?", telephone).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, result.Error
}
