package models

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	ID           uint       `gorm:"primaryKey"`
	DateTime     time.Time  `json:"dateTime"`
	TableNum     int        `json:"tableNum"`
	ExitTime     time.Time  `json:"exitTime"`
	UserID       uint       `json:"userId"`
	User         User       `gorm:"foreignKey:UserID" json:"user"`
	RestaurantID uint       `json:"restaurantId"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant"`
	gorm.Model   `json:"-" swaggerignore:"true"`
}

type ReservationHandler struct {
	db *gorm.DB
}

func NewReservationHandler(db *gorm.DB) *ReservationHandler {
	return &ReservationHandler{db}
}

func (h *ReservationHandler) CreateReservation(userID uint, reservation *Reservation) error {
	reservation.UserID = userID

	if err := h.db.Create(reservation).Error; err != nil {
		return err
	}

	return h.db.Preload("User").Preload("Restaurant").First(reservation, reservation.ID).Error
}

func (h *ReservationHandler) GetReservation(id uint) (*Reservation, error) {
	var reservation Reservation
	result := h.db.Preload("User").Preload("Restaurant").First(&reservation, id)
	return &reservation, result.Error
}

func (h *ReservationHandler) GetReservations() ([]Reservation, error) {
	var reservations []Reservation
	result := h.db.Preload("User").Preload("Restaurant").Find(&reservations)
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

func (handler *ReservationHandler) GetReservationsByUserID(userID uint) ([]Reservation, error) {
	var reservations []Reservation
	result := handler.db.Preload("User").Preload("Restaurant").Where("user_id = ?", userID).Find(&reservations)

	if result.Error != nil {
		return nil, result.Error
	}
	return reservations, nil
}
