package middleware

import (
	"net/http"

	authzservice "github.com/LavaJover/shvark-api-gateway/internal/client/authz-service"
	"github.com/gin-gonic/gin"
)

func RequirePermission(authzClient *authzservice.AuthzClient, object string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDAny, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID, ok := userIDAny.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid userID context"})
			return
		}

		resp, err := authzClient.CheckPermission(userID, object, action)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "authz error"})
			return
		}

		if !resp.Allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		c.Next()
	}
}