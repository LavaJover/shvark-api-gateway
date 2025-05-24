package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimitMiddleware() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit: 100,
	}
	store := memory.NewStore()
	limiter := limiter.New(store, rate)

	return func(c *gin.Context) {
		context, err := limiter.Get(c, c.ClientIP())
		if err != nil || context.Reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}
		c.Next()
	}
}