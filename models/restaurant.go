package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Restaurant struct {
	ID           uint     `gorm:"primaryKey"`
	Name         string   `json:"name"`
	Address      string   `json:"address"`
	Telephone    string   `json:"telephone"`
	OpenTime     string   `json:"openTime"`
	CloseTime    string   `json:"closeTime"`
	Instagram    string   `json:"instagram"`
	Facebook     string   `json:"facebook"`
	Description  string   `json:"description"`
	Rating       *float64 `json:"rating" gorm:"default:0" validate:"required,min=0"`
	CommentCount *float64 `json:"commentCount" gorm:"default:0" validate:"required,min=0"`
	ImageURL     string   `json:"imageUrl"`
	gorm.Model   `json:"-" swaggerignore:"true"`
}

type RestaurantHandler struct {
	db *gorm.DB
}

func NewRestaurantHandler(db *gorm.DB) *RestaurantHandler {
	return &RestaurantHandler{db}
}

func (h *RestaurantHandler) CreateRestaurant(restaurant *Restaurant) error {
	return h.db.Create(restaurant).Error
}

func (h *RestaurantHandler) GetRestaurant(id uint) (*Restaurant, error) {
	var restaurant Restaurant
	result := h.db.First(&restaurant, id)
	return &restaurant, result.Error
}

func (h *RestaurantHandler) GetRestaurants() ([]Restaurant, error) {
	var restaurants []Restaurant
	result := h.db.Find(&restaurants)
	return restaurants, result.Error
}

func (h *RestaurantHandler) UpdateRestaurant(id uint, restaurant *Restaurant) error {
	result := h.db.Model(&Restaurant{}).Where("id = ?", id).Updates(restaurant)
	return result.Error
}

func (h *RestaurantHandler) DeleteRestaurant(id uint) error {
	// Bypass soft delete and force a hard delete
	result := h.db.Unscoped().Where("id = ?", id).Delete(&Restaurant{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no restaurant found with id %d", id)
	}
	return nil
}
