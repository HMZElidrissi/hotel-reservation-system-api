package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Role(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "claims not found in context"})
			c.Abort()
			return
		}

		claimsMap, ok := claims.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "claims is not a map"})
			c.Abort()
			return
		}

		roleClaim, exists := claimsMap["role"]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role claim not found in claims"})
			c.Abort()
			return
		}

		if roleClaim != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "user does not have the required role"})
			c.Abort()
			return
		}

		c.Next()
	}
}
