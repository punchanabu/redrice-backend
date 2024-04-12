package models

import (
	"time"
	"gorm.io/gorm"
)

type Comment struct {
	ID           uint       `gorm:"primaryKey"`
	DateTime     time.Time  `json:"dateTime"`
	MyComment	 string 	`json:"myComment"`
	UserID       uint       `json:"userId"`
	User         User       `gorm:"foreignKey:UserID" json:"user"`
	RestaurantID uint       `json:"restaurantId"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant"`
	gorm.Model   `json:"-" swaggerignore:"true"`
}

type CommentHandler struct {
	db *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{db}
}

func (h *CommentHandler) CreateComment(userID uint, comment *Comment) error {
	comment.UserID = userID

	if err := h.db.Create(comment).Error; err != nil {
		return err
	}

	return h.db.Preload("User").Preload("Restaurant").First(comment, comment.ID).Error
}

func (h *CommentHandler) GetComment(id uint) (*Comment, error) {
	var comment Comment
	result := h.db.Preload("User").Preload("Restaurant").First(&comment, id)
	return &comment, result.Error
}

func (h *CommentHandler) GetComments() ([]Comment, error) {
	var comments []Comment
	result := h.db.Preload("User").Preload("Restaurant").Find(&comments)
	return comments, result.Error
}

func (h *CommentHandler) UpdateComment(id uint, comment *Comment) error {
	result := h.db.Model(&Comment{}).Where("id = ?", id).Updates(comment)
	return result.Error
}

func (h *CommentHandler) DeleteComment(id uint) error {
	result := h.db.Delete(&Comment{}, id)
	return result.Error
}

func (h *CommentHandler) GetCommentsByRestaurantID(restaurantID uint) ([]Comment, error) {
	var comments []Comment
	result := h.db.Preload("Restaurant").Preload("User").Where("restaurant_id = ?", restaurantID).Find(&comments)

	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}