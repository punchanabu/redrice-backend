package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/punchanabu/redrice-backend-go/middleware"
	"github.com/punchanabu/redrice-backend-go/models"
)

var userHandler *models.UserHandler

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

	if !userHandler.CheckPassword(user.Password, loginDetails.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	token, err := middleware.GenerateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(
		http.StatusOK, 
		gin.H{
			"token": token,
			"message" : "Login successful",
		},
	)
}
