package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/punchanabu/redrice-backend-go/models"
	"gorm.io/gorm"
)

var reservationHandler *models.ReservationHandler

func InitializedReservationHandler(db *gorm.DB) {
	reservationHandler = models.NewReservationHandler(db)
}

// @Summary Get a Single Reservation
// @Produce json
// @Param id path int true "Reservation ID"
// @Success 200 {object} models.Reservation
// @Failure 400 {object} string
// @Failure 404 {object} string
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

// @Summary Get a All Reservation
// @Produce json
// @Success 200 {object} []models.Reservation
// @Failure 500 {object} string
// @Router /reservations [get]
func GetReservations(c *gin.Context) {
	reservations, err := reservationHandler.GetReservations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching reservations!"})
		return
	}
	c.JSON(http.StatusOK, reservations)
}

// @Summary Create a Reservation
// @Produce json
// @Param reservation body models.Reservation true "Reservation"
// @Success 201 {object} models.Reservation
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /reservations [post]
func CreateReservation(c *gin.Context) {
	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User Id format"})
		return
	}

	err := reservationHandler.CreateReservation(uid, &reservation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating reservation"})
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

// @Summary Update a Reservation
// @Accept json
// @Produce json
// @Param id path int true "Reservation ID"
// @Param reservation body models.Reservation true "Reservation"
// @Success 200 {object} models.Reservation
// @Failure 400 {object} string
// @Failure 404 {object} string
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
// @Produce json
// @Param id path int true "Reservation ID"
// @Success 204 {object} string
// @Failure 400 {object} string
// @Failure 404 {object} string
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
