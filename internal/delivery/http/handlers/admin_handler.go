package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	adminRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/admin/request"
	adminResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/admin/response"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	SSOClient *client.SSOClient
	AuthzClient *client.AuthzClient
	OrderClient *client.OrderClient
	WalletClient *client.HTTPWalletClient
	UserClient *client.UserClient
}

func NewAdminHandler(
	ssoClient *client.SSOClient,
	authzClient *client.AuthzClient,
	orderClient *client.OrderClient,
	walletClient *client.HTTPWalletClient,
	userClient *client.UserClient,
) *AdminHandler {
	return &AdminHandler{
		SSOClient: ssoClient,
		AuthzClient: authzClient,
		OrderClient: orderClient,
		WalletClient: walletClient,
		UserClient: userClient,
	}
}

// @Summary Create new team
// @Description Create new team
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body adminRequest.CreateTeamRequest true "team credentials"
// @Success 201 {object} adminResponse.CreateTeamResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /admin/teams/create [post]
func (h *AdminHandler) CreateTeam(c *gin.Context) {
	var request adminRequest.CreateTeamRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	// register in sso
	registerResponse, err := h.SSOClient.Register(
		request.Login,
		request.Username,
		request.Password,
		"TRADER",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	// login to get access token
	loginResponse, err := h.SSOClient.Login(request.Login, request.Password, "")
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}
	// create wallet for new trader
	walletAddress, err := h.WalletClient.CreateWallet(registerResponse.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, adminResponse.CreateTeamResponse{
		TraderID: registerResponse.UserId,
		AccessToken: loginResponse.AccessToken,
		WalletAddress: walletAddress,
	})
}

// @Summary Create merchant account
// @Description Create merchant account
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body adminRequest.CreateMerchantRequest true "merchant credentials"
// @Success 201 {object} adminResponse.CreateMerchantResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /admin/merchants/create [post]
func (h *AdminHandler) CreateMerchant(c *gin.Context) {
	var request adminRequest.CreateMerchantRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	// register in sso
	registerResponse, err := h.SSOClient.Register(
		request.Login,
		request.Username,
		request.Password,
		"MERCHANT",
	)
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}
	loginResponse, err := h.SSOClient.Login(request.Login, request.Password, "")
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}
	// create wallet for merchant
	walletAddress, err := h.WalletClient.CreateWallet(registerResponse.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, adminResponse.CreateMerchantResponse{
		MerchantID: registerResponse.UserId,
		AccessToken: loginResponse.AccessToken,
		WalletAddress: walletAddress,
	})
}

// @Summary Create new traffic
// @Drescription Traffic is a relation between Merchant and Trader
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body adminRequest.CreateTrafficRequest true "create new traffic"
// @Success 201 {object} adminResponse.CreateTrafficResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /admin/traffic/create [post]
func (h *AdminHandler) CreateTraffic(c *gin.Context) {
	var request adminRequest.CreateTrafficRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	err := h.OrderClient.AddTraffic(
		request.MerchantID,
		request.TraderID,
		request.TraderReward,
		request.TraderPriority,
		request.PlatformFee,
		request.Enabled,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, adminResponse.CreateTrafficResponse{
		Message: "successfully created traffic",
	})
}

// @Summary Edit traffic
// @Description Edit traffic
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body adminRequest.EditTrafficRequest true "edit traffic"
// @Success 200 {object} adminResponse.EditTrafficResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /admin/traffic/edit [patch]
func (h *AdminHandler) EditTraffic(c *gin.Context) {
	var request adminRequest.EditTrafficRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	err := h.OrderClient.EditTraffic(
		request.Traffic.ID,
		request.Traffic.TraderReward,
		request.Traffic.TraderPriority,
		request.Traffic.PlatformFee,
		request.Traffic.Enabled,
	)
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, adminResponse.EditTrafficResponse{
	})
}

// @Summary Get traffic records
// @Description Get traffic records
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Success 200 {object} adminResponse.GetTrafficResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /admin/traffic/records [get]
func (h *AdminHandler) GetTrafficRecords(c *gin.Context) {
	var request adminRequest.GetTrafficRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if request.Page == 0 {
		request.Page = 1
	}
	if request.Limit == 0 {
		request.Limit = 1000
	}
	trafficResponse, err := h.OrderClient.GetTrafficRecords(
		request.Page,
		request.Limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	trafficRecords := make([]adminResponse.Traffic, len(trafficResponse))
	for i, trafficResp := range trafficResponse {
		trafficRecords[i] = adminResponse.Traffic{
			ID: trafficResp.Id,
			MerchantID: trafficResp.MerchantId,
			TraderID: trafficResp.TraderId,
			TraderReward: trafficResp.TraderRewardPercent,
			TraderPriority: trafficResp.TraderPriority,
			PlatformFee: trafficResp.PlatformFee,
			Enabled: trafficResp.Enabled,
		}
	}

	c.JSON(http.StatusOK, adminResponse.GetTrafficResponse{
		TrafficRecords: trafficRecords,
	})
}

// @Summary Create dispute
// @Description Open new dispute
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body adminRequest.CreateDisputeRequest true "new dispute data"
// @Success 201 {object} adminResponse.CreateDisputeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /admin/disputes/create [post]
func (h *AdminHandler) CreateDispute(c *gin.Context) {
	var request adminRequest.CreateDisputeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	disputeTtl, err := time.ParseDuration(request.Ttl)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	disputeID, err := h.OrderClient.CreateDispute(
		request.OrderID,
		request.ProofUrl,
		request.DisputeReason,
		disputeTtl,
		request.DisputeAmountFiat,
	)
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, adminResponse.CreateDisputeResponse{
		DisputeID: disputeID,
	})
}

// @Summary Accept active dispute
// @Description Accept active dispute
// @Tags admin
// @Security BearerAuth
// @Param input body adminRequest.AcceptDisputeRequest true "accept active dispute by ID"
// @Success 200 {object} adminResponse.AcceptDisputeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /admin/disputes/accept [post]
func (h *AdminHandler) AcceptDispute(c *gin.Context) {
	var request adminRequest.AcceptDisputeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	err := h.OrderClient.AcceptDispute(request.DisputeID)
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, adminResponse.AcceptDisputeResponse{})
}

// @Summary Reject active dispute
// @Description Reject active dispute
// @Tags admin
// @Security BearerAuth
// @Param input body adminRequest.RejectDisputeRequest true "reject active dispute by ID"
// @Success 200 {object} adminResponse.RejectDisputeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /admin/disputes/reject [post]
func (h *AdminHandler) RejectDispute(c *gin.Context) {
	var request adminRequest.RejectDisputeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	err := h.OrderClient.RejectDispute(request.DisputeID)
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, adminResponse.RejectDisputeResponse{})
}

// @Summary Get dispute info by dispute ID
// @Description Get dispute information by ID
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "dispute ID"
// @Success 200 {object} adminResponse.GetDisputeInfoResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /admin/disputes/{id} [get]
func (h *AdminHandler) GetDisputeInfo(c *gin.Context) {
	disputeID := c.Param("id")
	if disputeID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "empty path param ID"})
		return
	}
	disputeResponse, err := h.OrderClient.GetDisputeInfo(disputeID)
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, adminResponse.GetDisputeInfoResponse{
		Dispute: adminResponse.Dispute{
			DisputeID: disputeResponse.DisputeID,
			OrderID: disputeResponse.OrderID,
			ProofUrl: disputeResponse.ProofUrl,
			DisputeReason: disputeResponse.DisputeReason,
			DisputeStatus: disputeResponse.DisputeStatus,
		},
	})
}

// @Summary Freeze active dispute
// @Description Freeze active dispute
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body adminRequest.FreezeDisputeRequest true "dispute to Freeze"
// @Success 200 {object} adminResponse.FreezeDisputeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /admin/disputes/freeze [post]
func (h *AdminHandler) FreezeDispute(c *gin.Context) {
	var request adminRequest.FreezeDisputeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.OrderClient.FreeezeDispute(request.DisputeID)
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, adminResponse.FreezeDisputeResponse{})
}

// @Summary Get traders
// @Description Get traders
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} adminResponse.GetUsersResponse
// @Failure 502 {object} ErrorResponse
// @Router /admin/traders [get]
func (h *AdminHandler) GetTraders(c *gin.Context) {
	response, err := h.UserClient.GetTraders()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	tradersResp := make([]adminResponse.User, len(response.Traders))

	for i, trader := range response.Traders {
		tradersResp[i] = adminResponse.User{
			ID: trader.UserId,
			Username: trader.Username,
			Login: trader.Login,
			Role: trader.Role,
		}
	}

	c.JSON(http.StatusOK, adminResponse.GetUsersResponse{
		Users: tradersResp,
	})
}

// @Summary Get merchants
// @Description Get merchants
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} adminResponse.GetUsersResponse
// @Failure 502 {object} ErrorResponse
// @Router /admin/merchants [get]
func (h *AdminHandler) GetMerchants(c *gin.Context) {
	response, err := h.UserClient.GetMerchants()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	merchantsResp := make([]adminResponse.User, len(response.Merchants))

	for i, merchant := range response.Merchants {
		merchantsResp[i] = adminResponse.User{
			ID: merchant.UserId,
			Username: merchant.Username,
			Login: merchant.Login,
			Role: merchant.Role,
		}
	}

	c.JSON(http.StatusOK, adminResponse.GetUsersResponse{
		Users: merchantsResp,
	})
}

// @Summary Get order disputes
// @Descriptions get order disputes
// @Tags admin
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Param status query string false "Filter by status"
// @Success 200 {object} adminResponse.GetOrderDisputesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /admin/orders/disputes [get]
func (h *AdminHandler) GetOrderDisputes(c *gin.Context) {
	// Параметры по умолчанию
	const (
		defaultPage  = 1
		defaultLimit = 10
		maxLimit     = 100
	)

	page, err := strconv.Atoi(c.DefaultQuery("page", fmt.Sprintf("%d", defaultPage)))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid page parameter"})
		return
	}

	// Получаем и парсим limit
	limit, err := strconv.Atoi(c.DefaultQuery("limit", fmt.Sprintf("%d", defaultLimit)))
	if err != nil || limit < 1 || limit > maxLimit {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid limit parameter"})
		return
	}

	// Получаем статус (необязательный)
	status := c.Query("status")

	response, err := h.OrderClient.GetOrderDisputes(int64(page), int64(limit), status)
	if err != nil {
		c.JSON(http.StatusBadGateway,ErrorResponse{Error: err.Error()})
		return
	}

	disputes := make([]adminResponse.Dispute, len(response.Disputes))
	for i, dispute := range response.Disputes {
		disputes[i] = adminResponse.Dispute{
			DisputeID: dispute.DisputeId,
			ProofUrl: dispute.ProofUrl,
			DisputeReason: dispute.DisputeReason,
			DisputeStatus: dispute.DisputeStatus,
			DisputeAmountFiat: dispute.DisputeAmountFiat,
			DisputeAmountCrypto: dispute.DisputeAmountCrypto,
			DisputeCryptoRate: dispute.DisputeCryptoRate,
			OrderID: dispute.OrderId,
			AcceptAt: dispute.AcceptAt.AsTime(),
			Order: adminResponse.Order{
				ID: dispute.Order.OrderId,
				MerchantOrderID: dispute.Order.MerchantOrderId,
				AmountFiat: dispute.Order.AmountFiat,
				CryproRate: dispute.Order.CryptoRubRate,
				AmountCrypto: dispute.Order.AmountCrypto,
				BankDetail: adminResponse.BankDetail{
					BankName: dispute.Order.BankDetail.BankName,
					PaymentSystem: dispute.Order.BankDetail.PaymentSystem,
					Phone: dispute.Order.BankDetail.Phone,
					CardNumber: dispute.Order.BankDetail.CardNumber,
					Owner: dispute.Order.BankDetail.Owner,
					TraderID: dispute.Order.BankDetail.TraderId,
				},
			},
		}
	}

	c.JSON(http.StatusOK, adminResponse.GetOrderDisputesResponse{
		Disputes: disputes,
		Pagination: adminResponse.Pagination{
			CurrentPage: response.Pagination.CurrentPage,
			TotalPages: response.Pagination.TotalPages,
			TotalItems: response.Pagination.TotalItems,
			ItemsPerPage: response.Pagination.ItemsPerPage,
		},
	})
}

// @Summary Set withdrawal rules
// @Description Set withdrawal rules for user
// @Tags admin
// @Accept json
// @Produce json
// @Param input body adminRequest.SetWithdrawalRulesRequest true "Withdrawal rules"
// @Success 200 {object} adminResponse.SetWithdrawalRulesResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/wallets/withdraw/rules [post]
func (h *AdminHandler) SetWithdrawalRules(c *gin.Context) {
	var request adminRequest.SetWithdrawalRulesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	walletRequest := client.SetWithdrawalRulesRequest{
		TraderID: request.UserID,
		FixedFee: request.FixedFee,
		MinAmount: request.MinAmount,
		CooldownSeconds: request.CooldownSeconds,
	}

	resp, err := h.WalletClient.SetWithdrawalRules(&walletRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, adminResponse.SetWithdrawalRulesResponse{
		Success: resp.Success,
		Rule: adminResponse.Rule{
			ID: resp.Rule.ID,
			TraderID: resp.Rule.TraderID,
			FixedFee: resp.Rule.FixedFee,
			MinAmount: resp.Rule.MinAmount,
			CooldownSeconds: resp.Rule.CooldownSeconds,
			UpdatedAt: resp.Rule.UpdatedAt,
			CreatedAt: resp.Rule.CreatedAt,
		},
	})
}

// @Summary Get withdrawal rules
// @Description Get withdrawal rules for a given trader
// @Tags admin
// @Accept json
// @Produce json
// @Param userId path string true "user ID"
// @Success 200 {object} adminResponse.GetWithdrawalRulesResponse
// @Success 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/wallets/withdraw/rules/{userId} [get]
func (h *AdminHandler) GetUserWithdrawalRules(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID path param missed"})
		return
	}

	resp, err := h.WalletClient.GetWithdrawalRules(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, adminResponse.GetWithdrawalRulesResponse{
		Rule: adminResponse.Rule{
			ID: resp.ID,
			TraderID: resp.TraderID,
			FixedFee: resp.FixedFee,
			MinAmount: resp.MinAmount,
			CooldownSeconds: resp.CooldownSeconds,
			UpdatedAt: resp.UpdatedAt,
			CreatedAt: resp.CreatedAt,
		},
	})
}

// @Summary Delete withdrawal rule
// @Description Delete user withdrawal rule
// @Tags admin
// @Accept json
// @Produce json
// @Param userId path string true "user ID"
// @Success 200 {object} adminResponse.DeleteWithdrawalRulesResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/wallets/withdraw/rules/{userId} [delete]
func (h *AdminHandler) DeleteUserWithdrawalRules(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID path param missed"})
		return
	}

	err := h.WalletClient.DeleteWithdrawalRule(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, adminResponse.DeleteWithdrawalRulesResponse{})
}