package v1

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/punchanabu/redrice-backend-go/models"
	"gorm.io/gorm"
)

var commentHandler *models.CommentHandler
var restaurantHandler *models.RestaurantHandler

func InitializedCommentHandler(db *gorm.DB) {
	commentHandler = models.NewCommentHandler(db)
	restaurantHandler = models.NewRestaurantHandler(db)
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
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	log.Println("JSON bound successfully:", comment)

	userID, exist := c.Get("id")
	if !exist {
		log.Println("No user ID present in the request context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	log.Println("User ID found:", userID)

	uid, ok := userID.(uint)
	if !ok {
		log.Println("User ID is of invalid type:", userID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User Id format"})
		return
	}

	log.Println("User ID type assertion successful:", uid)

	err := commentHandler.CreateComment(uid, &comment)
	if err != nil {
		log.Println("Error creating comment:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating comment"})
		return
	}

	log.Println("Comment created successfully:", comment)

	// Assuming restaurantHandler is correctly instantiated and not nil
	restaurant, err := restaurantHandler.GetRestaurant(comment.RestaurantID)
	if err != nil {
		log.Println("Error fetching restaurant:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching restaurant for comment"})
		return
	}

	log.Println("Restaurant fetched successfully:", restaurant)

	// Update the rating and comment count
	restaurant.Rating = (restaurant.Rating*float64(restaurant.CommentCount) + comment.Rating) / float64(restaurant.CommentCount+1)
	restaurant.CommentCount++
	err = restaurantHandler.UpdateRestaurant(comment.RestaurantID, restaurant)
	if err != nil {
		log.Println("Error updating restaurant:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating restaurant comment count"})
		return
	}

	log.Println("Restaurant updated successfully:", restaurant)

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

	// Check if the user is the owner of the comment
	id, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}

	userID, ok := id.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot Parse User ID"})
	}

	ownComment, err := commentHandler.GetComment(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching comment for restaurant"})
		return
	}

	if ownComment.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to update this comment"})
		return
	}

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

	// Check if the user is the owner of the comment
	id, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}

	userID, ok := id.(uint)
	ownComment, err := commentHandler.GetComment(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching comment for restaurant"})
		return
	}

	if ownComment.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to delete this comment"})
		return
	}

	restaurant, err := restaurantHandler.GetRestaurant(ownComment.RestaurantID)
	if err != nil {
		log.Println("Error fetching restaurant:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching restaurant for comment"})
		return
	}

	log.Println("Restaurant fetched successfully:", restaurant)

	// Update the rating and comment count
	restaurant.Rating = (restaurant.Rating*float64(restaurant.CommentCount) - ownComment.Rating) / float64(restaurant.CommentCount-1)
	restaurant.CommentCount--
	err = restaurantHandler.UpdateRestaurant(ownComment.RestaurantID, restaurant)
	if err != nil {
		log.Println("Error updating restaurant:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating restaurant comment count"})
		return
	}

	err = commentHandler.DeleteComment(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// GetRestaurantComments retrieves all comments for a given restaurant ID.
// @Summary Get Reataurant's Comments
// @Description Retrieves a list of comments associated with a specific restaurant.
// @Tags comments
// @Produce json
// @Param restaurantId path int true "Reataurant ID"
// @Security Bearer
// @Success 200 {array} models.Comment "An array of comment objects for the restaurant."
// @Failure 400 {object} ErrorResponse "Invalid reataurant ID format."
// @Failure 404 {object} ErrorResponse "Comments not found for the specified restaurant ID."
// @Router /restaurants/{restaurantID}/comments [get]
func GetRestaurantComments(c *gin.Context) {
	RestaurantID := c.Param("id")
	if RestaurantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID format"})
		return
	}

	uid, err := strconv.ParseUint(RestaurantID, 10, 32) // Convert ReataurantID from string to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing restaurant ID"})
		return
	}

	comments, err := commentHandler.GetCommentsByRestaurantID(uint(uid)) // Correctly cast to uint now
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching comments for restaurant"})
		return
	}

	c.JSON(http.StatusOK, comments)
}
