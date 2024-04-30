package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/punchanabu/redrice-backend-go/middleware"
	"github.com/punchanabu/redrice-backend-go/models"
	"gorm.io/gorm"
)

var reservationHandler *models.ReservationHandler

func InitializedReservationHandler(db *gorm.DB) {
	reservationHandler = models.NewReservationHandler(db)
}

// @Summary Get a Single Reservation
// @Description Retrieves details of a single reservation by its unique identifier.
// @Tags reservations
// @Produce json
// @Param id path int true "Reservation ID" Format(int64)
// @security BearerAuth
// @Success 200 {object} models.Reservation "The details of the reservation including ID, DateTime, UserID, User, RestaurantID, and Restaurant."
// @Failure 400 {object} ErrorResponse "Invalid reservation ID format."
// @Failure 404 {object} ErrorResponse "Reservation not found with the specified ID."
// @Router /reservations/{id} [get]
func GetReservation(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation id"})
		return
	}

	idUint := uint(idInt)
	reservation, err := reservationHandler.GetReservation(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// @Summary Get All Reservations
// @Description Retrieves a list of all reservations in the system.
// @Tags reservations
// @Produce json
// @security BearerAuth
// @Success 200 {array} models.Reservation "An array of reservation objects."
// @Failure 500 {object} ErrorResponse "Internal server error while fetching reservations."
// @Router /reservations [get]
func GetReservations(c *gin.Context) {
	reservations, err := reservationHandler.GetReservations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching reservations!"})
		return
	}
	c.JSON(http.StatusOK, reservations)
}

// @Summary Create a New Reservation
// @Description Adds a new reservation to the system with the provided details. This endpoint requires authentication.
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation body models.Reservation true "Reservation Details"
// @security BearerAuth
// @Success 201 {object} models.Reservation "The created reservation's details, including its unique identifier."
// @Failure 400 {object} ErrorResponse "Invalid input format for reservation details."
// @Failure 500 {object} ErrorResponse "Internal server error while creating the reservation."
// @Router /reservations [post]
func CreateReservation(c *gin.Context) {
	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
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

	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	tokenString := parts[1]
	claims, err := middleware.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized 🥹 Please login first!"})
		c.Abort()
		return
	}

	OwnReservations, err := reservationHandler.GetReservationsByUserID(uint(uid)) // Correctly cast to uint now
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching reservations for user"})
		return
	}

	if len(OwnReservations) == 3 && claims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "User already has 3 reservations. Cannot create more."})
		return
	}

	err = reservationHandler.CreateReservation(uid, &reservation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating reservation"})
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

// @Summary Update a Reservation
// @Description Updates the details of an existing reservation identified by its ID. This endpoint requires authentication.
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path int true "Reservation ID" Format(int64)
// @Param reservation body models.Reservation true "Updated Reservation Details"
// @security BearerAuth
// @Success 200 {object} models.Reservation "The updated reservation's details."
// @Failure 400 {object} ErrorResponse "Invalid input format for reservation details or invalid reservation ID."
// @Failure 404 {object} ErrorResponse "Reservation not found with the specified ID."
// @Router /reservations/{id} [put]
func UpdateReservation(c *gin.Context) {
	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation id"})
		return
	}

	idUint := uint(idInt)

	err = reservationHandler.UpdateReservation(idUint, &reservation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating reservation"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// @Summary Delete a Reservation
// @Description Removes a reservation from the system by its unique identifier. This endpoint requires authentication.
// @Tags reservations
// @Produce json
// @Param id path int true "Reservation ID" Format(int64)
// @security BearerAuth
// @Success 204 "Reservation successfully deleted, no content to return."
// @Failure 400 {object} ErrorResponse "Invalid reservation ID format."
// @Failure 404 {object} ErrorResponse "Reservation not found with the specified ID."
// @Router /reservations/{id} [delete]
func DeleteReservation(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation id"})
		return
	}

	idUint := uint(idInt)

	err = reservationHandler.DeleteReservation(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation deleted successfully"})
}

// GetUserReservations retrieves all reservations for a given user ID.
// @Summary Get User's Reservations
// @Description Retrieves a list of reservations associated with a specific user.
// @Tags reservations
// @Produce json
// @Param userId path int true "User ID"
// @security BearerAuth
// @Success 200 {array} models.Reservation "An array of reservation objects for the user."
// @Failure 400 {object} ErrorResponse "Invalid user ID format."
// @Failure 404 {object} ErrorResponse "Reservations not found for the specified user ID."
// @Router /users/{userId}/reservations [get]
func GetUserReservations(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	uid, err := strconv.ParseUint(userID, 10, 32) // Convert userID from string to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing user ID"})
		return
	}

	reservations, err := reservationHandler.GetReservationsByUserID(uint(uid)) // Correctly cast to uint now
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching reservations for user"})
		return
	}

	c.JSON(http.StatusOK, reservations)
}
