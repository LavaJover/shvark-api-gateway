package handlers

import (
	"net/http"
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
}

func NewAdminHandler(
	ssoClient *client.SSOClient,
	authzClient *client.AuthzClient,
	orderClient *client.OrderClient,
	walletClient *client.HTTPWalletClient,
) *AdminHandler {
	return &AdminHandler{
		SSOClient: ssoClient,
		AuthzClient: authzClient,
		OrderClient: orderClient,
		WalletClient: walletClient,
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