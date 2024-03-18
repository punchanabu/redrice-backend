package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/punchanabu/redrice-backend-go/middleware"
	"github.com/punchanabu/redrice-backend-go/models"
	"gorm.io/gorm"
)

var userHandler *models.UserHandler

func InitializedAuthHandler(db *gorm.DB) {
	userHandler = models.NewUserHandler(db)
}

// @Summary Register a new user
// @Accept json
// @Produce json
// @Param user body models.User true "Register {name: string, telephone: string, email: string, password: string, role: string}"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /auth/register [post]
func Register(c *gin.Context) {

	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input format! please check the input format"})
		return
	}

	err := userHandler.CreateUser(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}
}

// @Summary Login a user
// @Description Login with email and password
// @Accept  json
// @Produce  json
// @Param user body object true "Login Credentials" { "email": "string", "password": "string" }
// @Success 200 {object} string "Login successful"
// @Failure 400 {object} string "Invalid input format"
// @Failure 401 {object} string "Authentication failed"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Server error"
// @Router /auth/login [post]
func Login(c *gin.Context) {

	var loginDetails struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input format! please check the input format"})
		return
	}

	user, err := userHandler.GetUserByEmail(loginDetails.Email)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Authentication failed"})
		return
	}

	if !userHandler.CheckPassword(user.Email, loginDetails.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password is incorrect!"})
		return
	}

	token, err := middleware.GenerateToken(user.Email, user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"token":   token,
			"message": "Login successful",
		},
	)
}
