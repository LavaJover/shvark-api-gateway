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
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	// init orders-client
	ordersAddr := "localhost:50058"
	ordersHandler, err := handlers.NewOrderHandler(ordersAddr)
	if err != nil {
		log.Printf("failed to init orders handler: %v\n", err)
	}

	// init wallet client
	walletHandler, err := handlers.NewWalletHandler()
	if err != nil {
		log.Printf("failed to init wallet client")
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

	// auth-service
	r.POST("/api/v1/register", authHandler.Register)
	r.POST("/api/v1/login", authHandler.Login)
	r.POST("/api/v1/validate_token", authHandler.ValidateToken)
	r.GET("/api/v1/user_by_token", authHandler.GetUserByToken)

	// user-service
	r.GET("/api/v1/users/:id", userHandler.GetUserByID)

	// profile-service
	r.GET("/api/v1/profiles/:uuid", profileHandler.GetProfileByID)

	// RBAC-service
	rbacGroup := r.Group("/api/v1/rbac", middleware.AuthMiddleware(authHandler.SSOClient))
	{
		rbacGroup.POST("/roles", authzHandler.AssignRole)
		rbacGroup.DELETE("/roles", authzHandler.RevokeRole)
		rbacGroup.POST("/policies", authzHandler.AddPolicy)
		rbacGroup.DELETE("/policies", authzHandler.DeletePolicy)
		rbacGroup.POST("/permissions", authzHandler.CheckPermission)
	}

	// banking-service
	bankingGroup := r.Group("/api/v1/banking", middleware.AuthMiddleware(authHandler.SSOClient))
	{
		bankingGroup.POST("/details", bankingHandler.CreateBankDetail)
		bankingGroup.POST("/details/delete", bankingHandler.DeleteBankDetail)
		bankingGroup.GET("/details/:uuid", bankingHandler.GetBankDetailByID)
		bankingGroup.PATCH("/details", bankingHandler.UpdateBankDetail)
		bankingGroup.GET("/details", bankingHandler.GetBankDetailsByTraderID)	
	}

	// orders-service
	orderGroup := r.Group("/api/v1/orders", middleware.AuthMiddleware(authHandler.SSOClient))
	{
		orderGroup.POST("/", ordersHandler.CreateOrder)
		orderGroup.GET("/:uuid", ordersHandler.GetOrderByID)
		orderGroup.GET("/trader/:traderUUID", ordersHandler.GetOrdersByTraderID)
		orderGroup.POST("/approve", ordersHandler.ApproveOrder)
		orderGroup.POST("/cancel", ordersHandler.CancelOrder)
		orderGroup.POST("/disputes/open", ordersHandler.OpenOrderDispute)
		orderGroup.POST("/disputes/resolve", ordersHandler.ResolveOrderDispute)	
	}

	// wallet-service
	walletGroup := r.Group("/api/v1/wallets", middleware.AuthMiddleware(authHandler.SSOClient))
	{
		walletGroup.POST("/create", walletHandler.CreateWallet)
		walletGroup.POST("/freeze", walletHandler.Freeze)
		walletGroup.POST("/release", walletHandler.Release)
		walletGroup.POST("/withdraw", middleware.RequirePermission(authzHandler.AuthzClient, "wallet", "withdraw"), walletHandler.Withdraw)
		walletGroup.POST("/deposit", walletHandler.Deposit)
		walletGroup.GET("/:traderID/history", middleware.RequireSelfOrAdmin(authzHandler.AuthzClient, "traderID"), walletHandler.GetTraderHistory)
		walletGroup.GET("/:traderID/balance", middleware.RequireSelfOrAdmin(authzHandler.AuthzClient, "traderID"), walletHandler.GetTraderBalance)
		walletGroup.GET("/:traderID/address", middleware.RequireSelfOrAdmin(authzHandler.AuthzClient, "traderID"), walletHandler.GetTraderWalletAddress)
	}

	r.Run(":8080")
}