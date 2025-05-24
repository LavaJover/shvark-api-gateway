package main

import (
	"log"

	"github.com/LavaJover/shvark-api-gateway/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	ssoAddr := "localhost:50051"
	authHandler, err := handlers.NewAuthHandler(ssoAddr)
	if err != nil {
		log.Fatalf("failed to init auth handler: %v\n", err)
	}

	r := gin.Default()
	r.POST("/api/v1/register", authHandler.Register)
	r.POST("/api/v1/login", authHandler.Login)
	r.POST("/api/v1/validate_token", authHandler.ValidateToken)
	r.GET("/api/v1/user_by_token", authHandler.GetUserByToken)
	r.Run(":8080")
}