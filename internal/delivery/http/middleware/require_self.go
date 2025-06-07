package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireSelf(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDAny, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID, ok := userIDAny.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid context"})
			return
		}

		param := c.Param(paramName)
		if param != userID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		c.Next()
	}
}