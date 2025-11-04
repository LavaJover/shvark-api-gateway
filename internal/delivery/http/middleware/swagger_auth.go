package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func BasicAuth() gin.HandlerFunc {
    username := os.Getenv("SWAGGER_USER")
    password := os.Getenv("SWAGGER_PASS")

    return func(c *gin.Context) {
        // Разрешаем доступ к статическим файлам Swagger без аутентификации
        if c.Request.URL.Path == "/swagger/index.html" || 
           c.Request.URL.Path == "/swagger/swagger-ui.css" ||
           c.Request.URL.Path == "/swagger/swagger-ui-bundle.js" ||
           c.Request.URL.Path == "/swagger/swagger-ui-standalone-preset.js" ||
           c.Request.URL.Path == "/swagger/favicon-32x32.png" {
            c.Next()
            return
        }

        user, pass, ok := c.Request.BasicAuth()
        if !ok || 
            subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
            subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
            c.Header("WWW-Authenticate", `Basic realm="Swagger Docs"`)
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        c.Next()
    }
}