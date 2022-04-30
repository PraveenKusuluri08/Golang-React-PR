package helpers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func EndPoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		bearerToken := strings.SplitAfter(token, "Bearer")
		fmt.Println(bearerToken[1])
		actualToken := strings.Replace(bearerToken[1], " ", "", -1)
		if actualToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token not found"})
			c.Abort()
			return
		}
		claims, err := ValidateJwtAuthToken(actualToken)
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
