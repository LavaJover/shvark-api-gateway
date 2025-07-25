package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func BasicAuth() gin.HandlerFunc {
    username := os.Getenv("SWAGGER_USER") // Задайте в .env
    password := os.Getenv("SWAGGER_PASS")

    return func(c *gin.Context) {
        user, pass, ok := c.Request.BasicAuth()
        if !ok || 
            subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
            subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
            c.Header("WWW-Authenticate", `Basic realm="Swagger Docs"`)
            c.AbortWithStatus(http.StatusUnauthorized)
        }
        c.Next()
    }
}
