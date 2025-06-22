package handlers

import (
	"net/http"

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
// @Failure 500 {object} ErrorResponse
// @Router /admin/teams [post]
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
	loginResponse, err := h.SSOClient.Login(request.Login, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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

func (h *AdminHandler) CreateMerchant(c *gin.Context) {

}

func (h *AdminHandler) CreateTraffic(c *gin.Context) {

}

func (h *AdminHandler) EditTraffic(c *gin.Context) {

}