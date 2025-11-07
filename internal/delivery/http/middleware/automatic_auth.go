package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// automatic_auth_middleware.go
func AutomaticAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Printf("📱 [AUTOMATIC_AUTH] Processing automatic endpoint: %s", c.Request.URL.Path)
        
        // Для liveness и auth можем использовать упрощенную аутентификацию
        // или специальные токены устройств
        
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            // Проверяем другие возможные места где может быть токен
            authHeader = c.GetHeader("X-Device-Token")
            log.Printf("📱 [AUTOMATIC_AUTH] X-Device-Token: %s", authHeader)
        }
        
        // Логируем все заголовки для отладки
        for name, values := range c.Request.Header {
            if strings.Contains(strings.ToLower(name), "auth") || strings.Contains(strings.ToLower(name), "token") {
                log.Printf("📱 [AUTOMATIC_AUTH] Header %s: %v", name, values)
            }
        }
        
        // Временно разрешаем все запросы к liveness для тестирования
        if c.Request.URL.Path == "/api/v1/automatic/liveness" {
            log.Printf("📱 [AUTOMATIC_AUTH] Allowing liveness without auth for testing")
            c.Next()
            return
        }
        
        // Для остальных эндпоинтов используем обычную аутентификацию
        // Или верните 401 если нужна аутентификация
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "error": "Authentication required",
            "path":  c.Request.URL.Path,
        })
    }
}