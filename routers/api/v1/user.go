package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/punchanabu/redrice-backend-go/models"
	"gorm.io/gorm"
)

var userHandler *models.UserHandler

func InitializedUserHandler(db *gorm.DB) {
	userHandler = models.NewUserHandler(db)
}

// @Summary Get a Single User
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	idUint := uint(idInt)

	user, err := userHandler.GetUser(idUint)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Get a All User
// @Produce json
// @Success 200 {object} []models.User
// @Failure 500 {object} string
// @Router /users [get]
func GetUsers(c *gin.Context) {
	users, err := userHandler.GetUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users!"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Create a User
// @Accept json
// @Produce json
// @Param user body models.User true "User"
// @Success 201 {object} models.User
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := userHandler.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user!"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// @Summary Update a User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User"
// @Success 200 {object} models.User
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	idUint := uint(idInt)

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = userHandler.UpdateUser(idUint, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Delete a User
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 500 {object} string
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	idUint := uint(idInt)

	err = userHandler.DeleteUser(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}

	c.Status(http.StatusNoContent)
}
