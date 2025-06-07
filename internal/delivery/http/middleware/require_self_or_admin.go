package middleware

import (
	"net/http"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	"github.com/gin-gonic/gin"
)

func RequireSelfOrAdmin(authzClient *client.AuthzClient, paramName string) gin.HandlerFunc {
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
		if param == userID {
			c.Next()
			return
		}

		resp, err := authzClient.CheckPermission(userID, "wallets", "read:any")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Authorization service is unavailable now"})
			return
		}

		if !resp.Allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		c.Next()
	}
}