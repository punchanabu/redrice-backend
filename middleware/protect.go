package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
		    c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		    c.Abort()
		    return
		}
		
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
		    c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer <token>"})
		    c.Abort()
		    return
		}

		tokenString := parts[1]

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
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
		    c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		    c.Abort()
		    return
		}
		
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
		    c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer <token>"})
		    c.Abort()
		    return
		}

		tokenString := parts[1]

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