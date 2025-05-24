package main

import (
	"log"

	"github.com/LavaJover/shvark-api-gateway/internal/delivery/http/handlers"
	"github.com/LavaJover/shvark-api-gateway/internal/delivery/http/middleware"
	_ "github.com/LavaJover/shvark-api-gateway/pkg/docs"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Shvark API Gateway
// @version 1.0
// @description REST API for ShvarkPay
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// init clients
	ssoAddr := "localhost:50051"
	authHandler, err := handlers.NewAuthHandler(ssoAddr)
	if err != nil {
		log.Fatalf("failed to init auth handler: %v\n", err)
	}

	r := gin.Default()

	// use middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.LogginMiddleware())
	r.Use(middleware.RateLimitMiddleware())
	r.Use(middleware.HeaderCheckMiddleware())

	// define routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/v1/register", authHandler.Register)
	r.POST("/api/v1/login", authHandler.Login)
	r.POST("/api/v1/validate_token", authHandler.ValidateToken)
	r.GET("/api/v1/user_by_token", authHandler.GetUserByToken)

	r.Run(":8080")
}