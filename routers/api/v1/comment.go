package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/punchanabu/redrice-backend-go/models"
	"gorm.io/gorm"
)

var commentHandler *models.CommentHandler

func InitializedCommentHandler(db *gorm.DB) {
	commentHandler = models.NewCommentHandler(db)
}

// @Summary Get All Comments
// @Description Retrieves a list of all comments in the system.
// @Tags comments
// @Produce json
// @Security Bearer
// @Success 200 {array} models.Comment "An array of comment objects."
// @Failure 500 {object} ErrorResponse "Internal server error while fetching comments."
// @Router /comments [get]
func GetComments(c *gin.Context) {
	comments, err := commentHandler.GetComments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching comments!"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

// @Summary Get a Single Comment
// @Description Retrieves details of a single commnet by its unique identifier.
// @Tags comments
// @Produce json
// @Param id path int true "Comment ID" Format(int64)
// @Security Bearer
// @Success 200 {object} models.Comment "The details of the comment including ID, DateTime, Detail, UserID, User, RestaurantID, and Restaurant."
// @Failure 400 {object} ErrorResponse "Invalid comment ID format."
// @Failure 404 {object} ErrorResponse "Comment not found with the specified ID."
// @Router /comments/{id} [get]
func GetComment(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment id"})
		return
	}

	idUint := uint(idInt)
	comment, err := commentHandler.GetComment(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// @Summary Create a New Comment
// @Description Adds a new comment to the system with customer's opinion. This endpoint requires authentication.
// @Tags reservations
// @Accept json
// @Produce json
// @Param commnet body models.Comment true "Your Comment"
// @Security Bearer
// @Success 201 {object} models.Comment "The created comment's details, including its unique identifier."
// @Failure 400 {object} ErrorResponse "Invalid input format for reservation details."
// @Failure 500 {object} ErrorResponse "Internal server error while creating the reservation."
// @Router /comments [post]
func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	userID, exist := c.Get("id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User Id format"})
		return
	}

	err := commentHandler.CreateComment(uid, &comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// @Summary Update a Comment
// @Description Updates the details of an existing comment identified by its ID. This endpoint requires authentication.
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID" Format(int64)
// @Param comment body models.Comment true "Updated comment Details"
// @Security Bearer
// @Success 200 {object} models.Comment "The updated comment's details."
// @Failure 400 {object} ErrorResponse "Invalid input format for comment details or invalid comment ID."
// @Failure 404 {object} ErrorResponse "Comment not found with the specified ID."
// @Router /comments/{id} [put]
func UpdateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment id"})
		return
	}

	idUint := uint(idInt)

	err = commentHandler.UpdateComment(idUint, &comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// @Summary Delete a Comment
// @Description Removes a comment from the system. This endpoint requires authentication.
// @Tags comments
// @Produce json
// @Param id path int true "Comment ID" Format(int64)
// @Success 204 "Comment successfully deleted, no content to return."
// @Failure 400 {object} ErrorResponse "Invalid comment ID format."
// @Failure 404 {object} ErrorResponse "Comment not found with the specified ID."
// @Router /comments/{id} [delete]
func DeleteComment(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment id"})
		return
	}

	idUint := uint(idInt)

	err = commentHandler.DeleteComment(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}