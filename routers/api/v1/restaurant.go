package v1

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/punchanabu/redrice-backend-go/models"
	"gorm.io/gorm"
)

var RestaurantHandler *models.RestaurantHandler

func InitializedRestaurantHandler(db *gorm.DB) {
	RestaurantHandler = models.NewRestaurantHandler(db)
}

// @Summary Get a Single Restaurant
// @Description Retrieves details of a single restaurant by its unique identifier.
// @Tags restaurants
// @Produce json
// @Param id path int true "Restaurant ID" Format(int64)
// @Security Bearer
// @Success 200 {object} models.Restaurant "The details of the restaurant including ID, name, location, and other relevant information."
// @Failure 400 {object} ErrorResponse "Invalid restaurant ID format."
// @Failure 404 {object} ErrorResponse "Restaurant not found with the specified ID."
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

// @Summary Get All Restaurants
// @Description Retrieves a list of all restaurants in the system.
// @Tags restaurants
// @Produce json
// @Security Bearer
// @Success 200 {array} models.Restaurant "An array of restaurant objects."
// @Failure 500 {object} ErrorResponse "Internal server error while fetching restaurants."
// @Router /restaurants [get]
func GetRestaurants(c *gin.Context) {
	users, err := RestaurantHandler.GetRestaurants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching restaurants!"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// @Summary Create a New Restaurant
// @Description Adds a new restaurant to the system with the provided details.
// @Tags restaurants
// @Accept json
// @Produce json
// @Param restaurant body models.Restaurant true "Restaurant Registration Details"
// @Security Bearer
// @Success 201 {object} models.Restaurant "The created restaurant's details, including its unique identifier."
// @Failure 400 {object} ErrorResponse "Invalid input format for restaurant details."
// @Failure 500 {object} ErrorResponse "Internal server error while creating the restaurant."
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
// @Description Updates the details of an existing restaurant identified by its ID.
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path int true "Restaurant ID" Format(int64)
// @Param restaurant body models.Restaurant true "Updated Restaurant Details"
// @Security Bearer
// @Success 200 {object} models.Restaurant "The updated restaurant's details."
// @Failure 400 {object} ErrorResponse "Invalid input format for restaurant details or invalid restaurant ID."
// @Failure 404 {object} ErrorResponse "Restaurant not found with the specified ID."
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
// @Description Removes a restaurant from the system by its unique identifier.
// @Tags restaurants
// @Produce json
// @Param id path int true "Restaurant ID" Format(int64)
// @Security Bearer
// @Success 204 "Restaurant successfully deleted, no content to return."
// @Failure 400 {object} ErrorResponse "Invalid restaurant ID format."
// @Failure 404 {object} ErrorResponse "Restaurant not found with the specified ID."
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

	c.JSON(http.StatusOK, gin.H{"status": "deleted", "message": "Restaurant deleted successfully!"})
}
