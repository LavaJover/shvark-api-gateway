package handlers

import (
	"net/http"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	adminRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/admin/request"
	adminResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/admin/response"
	orderResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/order/response"
	orderpb "github.com/LavaJover/shvark-order-service/proto/gen"
	userpb "github.com/LavaJover/shvark-user-service/proto/gen"
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
// @Param traderId query string false "Filter by trader"
// @Param merchantId query string false "Filter by merchant"
// @Param disputeId query string false "Filter by dispute ID"
// @Param orderId query string false "Filter by order ID"
// @Success 200 {object} adminResponse.GetOrderDisputesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/orders/disputes [get]
func (h *AdminHandler) GetOrderDisputes(c *gin.Context) {
	var query adminRequest.GetOrderDisputesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := orderpb.GetOrderDisputesRequest{
		Page: int64(query.Page),
		Limit: int64(query.Limit),
		Status: query.Status,
		TraderId: query.TraderID,
		DisputeId: query.DisputeID,
		MerchantId: query.MerchantID,
		OrderId: query.OrderID,
	}
	response, err := h.OrderClient.GetOrderDisputes(&req)
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

// @Summary Create new team relation
// @Description Create new team relation
// @Tags admin
// @Accept json
// @Produce json
// @Param input body adminRequest.CreateTeamRelationRequest true "new relation"
// @Success 201 {string} string "Success"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/teams/relations/create [post] 
func (h *AdminHandler) CreateTeamRelation(c *gin.Context) {
	var request adminRequest.CreateTeamRelationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.OrderClient.CreateTeamRelation(
		&orderpb.CreateTeamRelationRequest{
			TeamLeadId: request.TeamLeadID,
			TraderId: request.TraderID,
			Commission: request.TeamRelationParams.Commission,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

// @Summary Update team relation
// @Description Update team relation
// @Tags admin
// @Accept json
// @Produce json
// @Param input body adminRequest.UpdateTeamRelationRequest true "update relation"
// @Success 200 {string} string "Success"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/teams/relations/update [patch]
func (h *AdminHandler) UpdateRelationParams(c *gin.Context) {
	var request adminRequest.UpdateTeamRelationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.OrderClient.UpdateTeamRelationParams(
		&orderpb.UpdateRelationParamsRequest{
			Relation: &orderpb.TeamRelationship{
				Id: request.RelationID,
				Commission: request.TeamRelationParams.Commission,
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// @Summary Get relations by teamLead id
// @Description get relations by teamLead ID
// @Tags admin
// @Accept json
// @Produce json
// @Param teamLeadID path string true "teamLeadID"
// @Success 200 {object} adminResponse.TeamRelationsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/teams/relations/team-lead/{teamLeadID} [get]
func (h *AdminHandler) GetRelationsByTeamLeadID(c *gin.Context) {
	teamLeadID := c.Param("teamLeadID")
	if teamLeadID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teamLeadID path param misses"})
		return
	}

	response, err := h.OrderClient.GetTeamRelationsByTeamLeadID(
		&orderpb.GetRelationsByTeamLeadIDRequest{
			TeamLeadId: teamLeadID,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	teamRelations := make([]adminResponse.TeamRelation, len(response.TeamRelations))
	for i, relation := range response.TeamRelations {
		teamRelations[i] = adminResponse.TeamRelation{
			ID: relation.Id,
			TraderID: relation.TraderId,
			TeamLeadID: relation.TeamLeadId,
			TeamRelationRarams: adminResponse.TeamRelationRarams{
				Commission: relation.Commission,
			},
		}
	}

	c.JSON(http.StatusOK, adminResponse.TeamRelationsResponse{
		TeamRelations: teamRelations,
	})
}

// @Summary Delete team relationship
// @Description Delete team relationship
// @Tags admin
// @Accept json
// @Produce json
// @Param relationID path string true "id of relationship"
// @Success 200 {string} string "Success"
// @Failure 404 {object} ErrorResponse
// @Router /admin/teams/relations/{relationID}/delete [delete]
func (h *AdminHandler) DeleteTeamRelationship(c *gin.Context) {
	relationID := c.Param("relationID")
	_, err := h.OrderClient.DeleteTeamRelationship(
		&orderpb.DeleteTeamRelationshipRequest{
			RelationId: relationID,
		},
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to delete relationship"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// @Summary Promote trader to team lead
// @Description Promote trader to team lead
// @Tags admin
// @Accept json
// @Produce json
// @Param traderID path string true "trader ID to be promoted to teamlead"
// @Success 200 {string} string "Success"
// @Failure 404 {object} ErrorResponse
// @Router /admin/teams/traders/{traderID}/promote-to-teamlead [post] 
func (h *AdminHandler) PromoteToTeamLead(c *gin.Context) {
	traderID := c.Param("traderID")
	_, err := h.UserClient.PromoteToTeamLead(
		&userpb.PromoteToTeamLeadRequest{
			UserId: traderID,
		},
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to promote trader to teamlead"})
		return
	}
	_, err = h.AuthzClient.AssignRole(traderID, "teamlead")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to assign teamlead role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// @Summary Demote teamlead
// @Description Demote teamlead to trader
// @Tags admin
// @Accept json
// @Produce json
// @Param teamleadID path string true "teamlead ID to be demoted"
// @Success 200 {string} string "Success"
// @Failure 404 {object} ErrorResponse
// @Router /admin/teams/teamleads/{teamleadID}/demote [post] 
func (h *AdminHandler) DemoteTeamLead(c *gin.Context) {
	teamleadID := c.Param("teamleadID")

	_, err := h.UserClient.DemoteTeamLead(
		&userpb.DemoteTeamLeadRequest{
			TeamLeadId: teamleadID,
		},
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to demote teamlead"})
		return
	}
	_, err = h.AuthzClient.RevokeRole(teamleadID, "teamlead")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to assign teamlead role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// @Summary Get users by role
// @Description Get users filtered by role
// @Tags admin
// @Accept json
// @Produce json
// @Param role query string false "user role"
// @Success 200 {object} adminResponse.GetUsersResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/users [get]
func (h *AdminHandler) GetUsersByRole(c *gin.Context) {
	role := c.Query("role")
	resp, err := h.UserClient.GetUsersByRole(
		&userpb.GetUsersByRoleRequest{
			Role: role,
		},
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to find users"})
		return
	}
	users := make([]adminResponse.User, len(resp.Users))
    for i, u := range resp.Users {
        users[i] = adminResponse.User{
            ID:        u.UserId,
            Login:     u.Login,
            Role:      u.Role,
			Username: u.Username,
        }
    }
    
    // Возвращаем успешный ответ
    c.JSON(http.StatusOK, adminResponse.GetUsersResponse{
        Users: users,
    })
}

// @Summary 	Get order statistics
// @Description Get order statistics
// @Tags 		admin
// @Security 	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		traderID  query string true "trader ID"
// @Param       date_from query string true "Дата начала (RFC3339 format, e.g. 2025-07-21T00:00:00Z)"
// @Param       date_to   query string true "Дата конца (RFC3339 format, e.g. 2025-07-21T23:59:59Z)"
// @Success 200 {object} orderResponse.GetOrderStatsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/orders/statistics [get]
func (h *AdminHandler) GetTraderOrderStats(c *gin.Context) {
	traderID := c.Query("traderID")

	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")

	dateFrom, err := time.Parse(time.RFC3339, dateFromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_from format, expected RFC3339"})
		return
	}

	dateTo, err := time.Parse(time.RFC3339, dateToStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_to format, expected RFC3339"})
		return
	}
	resp, err := h.OrderClient.GetOrderStats(
		traderID,
		dateFrom,
		dateTo,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no stats was found"})
		return
	}
	c.JSON(http.StatusOK, orderResponse.GetOrderStatsResponse{
		TotalOrders: resp.TotalOrders,
		SucceedOrders: resp.SucceedOrders,
		CanceledOrders: resp.CanceledOrders,
		ProcessedAmountFiat: float64(resp.ProcessedAmountFiat),
		ProcessedAmountCrypto: float64(resp.ProcessedAmountCrypto),
		CanceledAmountFiat: float64(resp.CanceledAmountFiat),
		CanceledAmountCrypto: float64(resp.CanceledAmountCrypto),
		IncomeCrypto: float64(resp.IncomeCrypto),
	})
}