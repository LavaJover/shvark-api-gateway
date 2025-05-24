package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LogginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// processing request
		c.Next()

		duration := time.Since(start)
		log.Printf(
            "[%s] %s %s %d %s",
            c.Request.Method,
            c.Request.URL.Path,
            c.ClientIP(),
            c.Writer.Status(),
            duration,
        )
	}
}