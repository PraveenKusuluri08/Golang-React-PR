package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func EndPoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token not found"})
			c.Abort()
			return
		}
		claims, err := ValidateJwtAuthToken(token)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})

			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
