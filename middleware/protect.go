package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized 🥹 Please login first!"})
			c.Abort()
			return
		}

		// Set user id to next handler for easy access
		c.Set("id", claims.UserId)
		c.Next()
	}
}

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized 🥹 Please login first!"})
			c.Abort()
			return
		}

		if claims.Role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized You are not an admin! 🥹 whahahahaha"})
			c.Abort()
			return
		}

		// Set user id to next handler for easy access
		c.Set("id", claims.UserId)
		c.Next()
	}
}