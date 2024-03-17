package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/punchanabu/redrice-backend-go/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var RestaurantHandler *models.RestaurantHandler

func InitializedRestaurantHandler(db *gorm.DB) {
	RestaurantHandler = models.NewRestaurantHandler(db)
}

// @Summary Get a Single Restaurant
// @Produce json
// @Param id path int true "Restaurant ID"
// @Success 200 {object} models.Restaurant
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Router /restaurants/{id} [get]
func GetRestaurant(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant id"})
		return
	}

	idUint := uint(idInt)
	restaurant, err := RestaurantHandler.GetRestaurant(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

// @Summary Get a All Restaurant
// @Produce json
// @Success 200 {object} []models.Restaurant
// @Failure 500 {object} string
// @Router /restaurants [get]
func GetRestaurants(c *gin.Context) {
	users, err := RestaurantHandler.GetRestaurants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching restaurants!"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// @Summary Create a Restaurant
// @Produce json
// @Param restaurant body models.Restaurant true "Restaurant"
// @Success 201 {object} models.Restaurant
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /restaurants [post]
func CreateRestaurant(c *gin.Context) {
	var restaurant models.Restaurant
	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := RestaurantHandler.CreateRestaurant(&restaurant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating restaurant!"})
		return
	}

	c.JSON(http.StatusCreated, restaurant)
}

// @Summary Update a Restaurant
// @Produce json
// @Param id path int true "Restaurant ID"
// @Param restaurant body models.Restaurant true "Restaurant"
// @Success 200 {object} models.Restaurant
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Router /restaurants/{id} [put]
func UpdateRestaurant(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant id"})
		return
	}

	idUint := uint(idInt)

	var restaurant models.Restaurant
	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = RestaurantHandler.UpdateRestaurant(idUint, &restaurant)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

// @Summary Delete a Restaurant
// @Produce json
// @Param id path int true "Restaurant ID"
// @Success 204 {object} string
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Router /restaurants/{id} [delete]
func DeleteRestaurant(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant id"})
		return
	}

	idUint := uint(idInt)

	err = RestaurantHandler.DeleteRestaurant(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"status": "deleted"})
}
