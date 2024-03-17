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

type ReservationHandler struct {
	db *gorm.DB
}

func NewReservationHandler(db *gorm.DB) *ReservationHandler {
	return &ReservationHandler{db}
}

func (h *ReservationHandler) CreateReservation(reservation *Reservation) error {
	return h.db.Create(reservation).Error
}

func (h *ReservationHandler) GetReservation(id uint) (*Reservation, error) {
	var reservation Reservation
	result := h.db.First(&reservation, id)
	return &reservation, result.Error
}

func (h *ReservationHandler) GetReservations() ([]Reservation, error) {
	var reservations []Reservation
	result := h.db.Find(&reservations)
	return reservations, result.Error
}

func (h *ReservationHandler) UpdateReservation(id uint, reservation *Reservation) error {
	result := h.db.Model(&Reservation{}).Where("id = ?", id).Updates(reservation)
	return result.Error
}

func (h *ReservationHandler) DeleteReservation(id uint) error {
	result := h.db.Delete(&Reservation{}, id)
	return result.Error
}
