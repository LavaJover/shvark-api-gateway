package middleware

import "github.com/gin-gonic/gin"

func HeaderCheckMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.GetHeader("Content-Type") != "application/json" {
            c.AbortWithStatusJSON(400, gin.H{"error": "invalid content type"})
            return
        }
        c.Next()
    }
}