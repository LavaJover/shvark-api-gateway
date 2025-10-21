package main

import (
	"fmt"
	"log"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	"github.com/LavaJover/shvark-api-gateway/internal/config"
	"github.com/LavaJover/shvark-api-gateway/internal/delivery/http/handlers"
	"github.com/LavaJover/shvark-api-gateway/internal/delivery/http/middleware"
	"github.com/LavaJover/shvark-api-gateway/pkg/docs"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 						Shvark API Gateway
// @version 					1.0
// @description 				REST API for ShvarkPay
// @host 						http://localhost:8080
// @schemes 					http
// @BasePath 					/api/v1
//
// @securityDefinitions.apikey 	BearerAuth
// @in 							header
// @name 						Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("failed to load .env")
	}

	cfg := config.MustLoad()

	// setup swagger based on development environment
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s",  cfg.SwaggerConfig.Host, cfg.SwaggerConfig.Port)
	docs.SwaggerInfo.Schemes = []string{cfg.SwaggerConfig.Schemes, "https"}
	docs.SwaggerInfo.BasePath = cfg.SwaggerConfig.BasePath

	// init user-client
	userAddr := fmt.Sprintf("%s:%s", cfg.UserService.Host, cfg.UserService.Port)
	userHandler, err := handlers.NewUserHandler(userAddr)
	if err != nil {
		log.Printf("failed to init user handler: %v", err)
	}

	// init sso-client
	ssoAddr := fmt.Sprintf("%s:%s", cfg.SSOService.Host, cfg.SSOService.Port)
	authHandler, err := handlers.NewAuthHandler(ssoAddr, userAddr)
	if err != nil {
		log.Printf("failed to init auth handler: %v\n", err)
	}

	// init authz-client
	authzAddr := fmt.Sprintf("%s:%s", cfg.AuthzService.Host, cfg.AuthzService.Port)
	authzHandler, err := handlers.NewAuthzhandler(authzAddr)
	if err != nil {
		log.Printf("failed to init authz handler: %v", err)
	}

	// init orders-client
	ordersAddr := fmt.Sprintf("%s:%s", cfg.OrderService.Host, cfg.OrderService.Port)
	ordersHandler, err := handlers.NewOrderHandler(ordersAddr)
	if err != nil {
		log.Printf("failed to init orders handler: %v\n", err)
	}
	
	bankingHandler, err := handlers.NewBankingHandler(ordersAddr)
	if err != nil {
		log.Printf("failed to init banking handler")
	}

	// init wallet client
	walletHandler, err := handlers.NewWalletHandler(client.NewHTTPWalletClient(fmt.Sprintf("%s:%s", cfg.WalletService.Host, cfg.WalletService.Port)))
	if err != nil {
		log.Printf("failed to init wallet client")
	}

	// init payments handlet
	paymentHandler, err := handlers.NewPaymentHandler(ordersHandler.OrderClient)
	if err != nil {
		log.Printf("failed to init payment handler")
	}

	r := gin.Default()

	// use middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.LogginMiddleware())
	r.Use(middleware.RateLimitMiddleware())
	// r.Use(middleware.HeaderCheckMiddleware())

	// define routes
	r.GET("/swagger/*any", middleware.BasicAuth(), ginSwagger.WrapHandler(swaggerFiles.Handler))

	// auth-service
	authGroup := r.Group("/api/v1")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/validate_token", authHandler.ValidateToken)
		authGroup.POST("/2fa/setup", middleware.AuthMiddleware(authHandler.SSOClient), authHandler.Setup2FA)
		authGroup.POST("/2fa/verify", middleware.AuthMiddleware(authHandler.SSOClient), authHandler.Verify2FA)
	}

	// user-service
	r.GET("/api/v1/users/:id", userHandler.GetUserByID)

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
		bankingGroup.GET("/details/stats/:traderID", bankingHandler.GetBankDetailsStats)
		bankingGroup.GET("/requisites", bankingHandler.GetBankDetails)
	}

	// orders-service
	orderGroup := r.Group("/api/v1/orders", middleware.AuthMiddleware(authHandler.SSOClient))
	{
		orderGroup.POST("/", ordersHandler.CreateOrder)
		orderGroup.GET("/:uuid", ordersHandler.GetOrderByID)
		orderGroup.GET("/trader/:traderUUID", ordersHandler.GetOrdersByTraderID)
		orderGroup.POST("/approve", ordersHandler.ApproveOrder)
		orderGroup.POST("/cancel", ordersHandler.CancelOrder)
		orderGroup.GET("/merchant/:id", ordersHandler.GetOrderByMerchantOrderID)
		orderGroup.GET("/statistics", ordersHandler.GetOrderStats)
		orderGroup.GET("/all", ordersHandler.GetAllOrders)
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
		walletGroup.POST("/offchain-withdraw", middleware.AuthMiddleware(authHandler.SSOClient), walletHandler.OffchainWithdraw)
		walletGroup.GET("/:traderID/commission-profit", walletHandler.GetCommissionProfit)
	}

	// payments for merchant
	paymentsGroup := r.Group("/api/v1/payments")
	{
		paymentsGroup.POST("/in/h2h", paymentHandler.CreateH2HPayIn)
		paymentsGroup.GET("/in/h2h/:id", paymentHandler.GetH2HPayInInfo)
		paymentsGroup.POST("/in/h2h/:id/cancel", paymentHandler.CancelPayIn)
		paymentsGroup.POST("/in/h2h/:id/arbitrage/link", paymentHandler.OpenPayInArbitrageLink)
		paymentsGroup.GET("/in/h2h/:id/arbitrage/info", paymentHandler.GetPayInArbitrageInfo)
	}

	walletAddr := fmt.Sprintf("%s:%s", cfg.WalletService.Host, cfg.WalletService.Port)
	walletClient := client.NewHTTPWalletClient(walletAddr)

	adminHandler := handlers.NewAdminHandler(
		authHandler.SSOClient,
		authzHandler.AuthzClient,
		ordersHandler.OrderClient,
		walletClient,
		userHandler.UserClient,
	)
	adminGroup := r.Group("/api/v1/admin")
	{
		adminGroup.POST("/teams/create", adminHandler.CreateTeam)
		adminGroup.POST("/merchants/create", adminHandler.CreateMerchant)
		adminGroup.POST("/traffic/create", adminHandler.CreateTraffic)
		adminGroup.PATCH("/traffic/edit", adminHandler.EditTraffic)
		adminGroup.DELETE("/traffic/:trafficId", adminHandler.DeleteTraffic)
		adminGroup.GET("/traffic/records", adminHandler.GetTrafficRecords)
		adminGroup.POST("/disputes/create", adminHandler.CreateDispute)
		adminGroup.POST("/disputes/accept", adminHandler.AcceptDispute)
		adminGroup.POST("/disputes/reject", adminHandler.RejectDispute)
		adminGroup.GET("/disputes/:id", adminHandler.GetDisputeInfo)
		adminGroup.POST("/disputes/freeze", adminHandler.FreezeDispute)
		adminGroup.GET("/traders", adminHandler.GetTraders)
		adminGroup.GET("/merchants", adminHandler.GetMerchants)
		adminGroup.GET("/orders/disputes", adminHandler.GetOrderDisputes)
		adminGroup.POST("/wallets/withdraw/rules", adminHandler.SetWithdrawalRules)
		adminGroup.GET("/wallets/withdraw/rules/:userId", adminHandler.GetUserWithdrawalRules)
		adminGroup.DELETE("/wallets/withdraw/rules/:userId", adminHandler.DeleteUserWithdrawalRules)
		adminGroup.POST("/teams/relations/create", adminHandler.CreateTeamRelation)
		adminGroup.PATCH("/teams/relations/update", adminHandler.UpdateRelationParams)
		adminGroup.GET("/teams/relations/team-lead/:teamLeadID", adminHandler.GetRelationsByTeamLeadID)
		adminGroup.DELETE("/teams/relations/:relationID/delete", adminHandler.DeleteTeamRelationship)
		adminGroup.POST("/teams/traders/:traderID/promote-to-teamlead", adminHandler.PromoteToTeamLead)
		adminGroup.POST("/teams/teamleads/:teamleadID/demote", adminHandler.DemoteTeamLead)
		adminGroup.GET("/users", adminHandler.GetUsersByRole)
		adminGroup.GET("/orders/statistics", adminHandler.GetTraderOrderStats)
	}

	merchantHandler := handlers.NewMerchanHandler(ordersHandler.OrderClient, walletClient, userHandler.UserClient, authHandler.SSOClient)
	merchantGroup := r.Group("/api/v1/merchant")
	{
		merchantGroup.POST("/order/:accountID/deposit", middleware.AuthMiddleware(authHandler.SSOClient), merchantHandler.CreatePayIn)
		merchantGroup.GET("/accounts/balance", middleware.AuthMiddleware(authHandler.SSOClient), merchantHandler.GetAccountBalance)
		merchantGroup.POST("/accounts/withdraw/create", middleware.AuthMiddleware(authHandler.SSOClient), merchantHandler.Withdraw)
		merchantGroup.GET("/banks", merchantHandler.GetBanks)
		merchantGroup.GET("/order/:iternalId/status", middleware.AuthMiddleware(authHandler.SSOClient), merchantHandler.GetOrderStatus)
		merchantGroup.POST("/auth/sign-in", merchantHandler.Login)
		merchantGroup.GET("/order", middleware.AuthMiddleware(authHandler.SSOClient), merchantHandler.GetOrders)
	}

	// init device handler
	deviceHandler, err := handlers.NewDeviceHandler(ordersHandler.OrderClient)
	if err != nil {
		log.Printf("failed to init device handler")
	}
	deviceGroup := r.Group("/api/v1/devices")
	{
		deviceGroup.POST("", deviceHandler.CreateDevice)
		deviceGroup.GET("/:traderId", deviceHandler.GetTraderDevices)
		deviceGroup.PATCH("/:deviceId/edit", deviceHandler.EditDevice)
		deviceGroup.DELETE("/:deviceId", deviceHandler.DeleteDevice)
	}

	automaticHandler := handlers.NewAutomaticHandler(adminHandler.OrderClient)
	automaticGroup := r.Group("/api/v1/automatic")
	{
		automaticGroup.POST("/process-sms", automaticHandler.Sms)
		automaticGroup.POST("/liveness", automaticHandler.Live)
		automaticGroup.POST("/auth", automaticHandler.Auth)
	}

	trafficHandler := handlers.NewTrafficHandler(adminHandler.OrderClient)
	trafficGroup := r.Group("/api/v1/traffic")
	{
		trafficGroup.PATCH("/traders/:traderID", trafficHandler.SetTraderLockTrafficStatus)
		trafficGroup.PATCH("/merchants/:merchantID", trafficHandler.SetMerchantLockTrafficStatus)
		trafficGroup.PATCH("/:trafficID/manual", trafficHandler.SetManuallyLockTrafficStatus)
		trafficGroup.PATCH("/antifraud/:traderID", trafficHandler.SetAntifraudLockTrafficStatus)
		trafficGroup.GET("/:trafficID/lock-statuses", trafficHandler.GetTrafficLockStatuses)
		trafficGroup.GET("/:trafficID/unlocked", trafficHandler.CheckTrafficUnlocked)
	}
	r.Run(":8080")
}