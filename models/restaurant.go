package models

import (
	"gorm.io/gorm"
)

type Restaurant struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	Description string `json:"description"`
	gorm.Model `json:"-" swaggerignore:"true"`
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
	result := h.db.Delete(&Restaurant{}, id)
	return result.Error
}

