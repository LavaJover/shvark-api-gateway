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
	// init sso-client
	ssoAddr := "localhost:50051"
	authHandler, err := handlers.NewAuthHandler(ssoAddr)
	if err != nil {
		log.Printf("failed to init auth handler: %v\n", err)
	}

	// init user-client
	userAddr := "localhost:50052"
	userHandler, err := handlers.NewUserHandler(userAddr)
	if err != nil {
		log.Printf("failed to init user handler")
	}

	// init profile-client
	profileAddr := "localhost:50055"
	profileHandler, err := handlers.NewProfileHandler(profileAddr)
	if err != nil {
		log.Printf("failed to init profile handler")
	}

	// init authz-client
	authzAddr := "localhost:50054"
	authzHandler, err := handlers.NewAuthzhandler(authzAddr)
	if err != nil {
		log.Printf("failed to init authz handler")
	}

	// init banking-client
	bankingAddr := "localhost:50057"
	bankingHandler, err := handlers.NewBankingHandler(bankingAddr)
	if err != nil {
		log.Printf("failed to init banking handler")
	}

	r := gin.Default()

	// use middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.LogginMiddleware())
	r.Use(middleware.RateLimitMiddleware())
	// r.Use(middleware.HeaderCheckMiddleware())

	// define routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/v1/register", authHandler.Register)
	r.POST("/api/v1/login", authHandler.Login)
	r.POST("/api/v1/validate_token", authHandler.ValidateToken)
	r.GET("/api/v1/user_by_token", authHandler.GetUserByToken)

	// user-service
	r.GET("/api/v1/users/:id", userHandler.GetUserByID)

	// profile-service
	r.GET("/api/v1/profiles/:uuid", profileHandler.GetProfileByID)

	// RBAC-service
	r.POST("/api/v1/rbac/roles", authzHandler.AssignRole)
	r.DELETE("/api/v1/rbac/roles", authzHandler.RevokeRole)
	r.POST("/api/v1/rbac/policies", authzHandler.AddPolicy)
	r.DELETE("/api/v1/rbac/policies", authzHandler.DeletePolicy)
	r.POST("/api/v1/rbac/permissions", authzHandler.CheckPermission)

	// banking-service
	r.POST("/api/v1/banking/details", bankingHandler.CreateBankDetail)
	r.GET("/api/v1/banking/details/:uuid", bankingHandler.GetBankDetailByID)

	r.Run(":8080")
}