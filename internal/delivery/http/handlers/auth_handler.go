package handlers

/*

To-Do:
1) gRPC microservice errors -> HTTP status
2) Retry middleware
3) Connetion pools - grpc.WithResolvers + balancing
4) OpenTelemetry for tracing
5) Metrics (callback time)
6) Viper configs
*/

import (
	"net/http"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	authRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/auth/request"
	authResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/auth/response"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	SSOClient *client.SSOClient
}

func NewAuthHandler(addr string) (*AuthHandler, error) {
	ssoClient, err := client.NewSSOClient(addr)
	if err != nil {
		return nil, err
	}

	return &AuthHandler{SSOClient: ssoClient}, nil
}

// @Summary User registration
// @Description Creating new account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body authRequest.RegisterRequest true "Registration credentials"
// @Success 201 {object} authResponse.RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var request authRequest.RegisterRequest

	// validate incoming request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// calling gRPC sso-service
	response, err := h.SSOClient.Register(request.Login, request.Username, request.Password, request.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, authResponse.RegisterResponse{
		UserID: response.UserId,
	})
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid data"`
}

// @Summary User login
// @Description Log in user account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body authRequest.LoginRequest true "Login credentials"
// @Success 200 {object} authResponse.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var request authRequest.LoginRequest

	// validate incoming request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// calling gRPC sso-service Login handler
	response, err := h.SSOClient.Login(request.Login, request.Password, request.TwoFACode)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code(){
			case codes.NotFound:
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "user not found",
				})
			}
		}
		if err.Error() == "rpc error: code = Unauthenticated desc = 2FA required" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "2FA_required",
			})
			return
		}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}

	c.JSON(http.StatusOK, authResponse.LoginResponse{
		AccessToken: response.AccessToken,
		RefreshToken: response.RefreshToken,
	})
}

// @Summary JWT validation check-point
// @Description Check if JWT is valid or not
// @Tags auth
// @Accept json
// @Produce json
// @Param input body authRequest.ValidateTokenRequest true "JWT instance"
// @Success 200 {object} authResponse.ValidateTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /validate_token [post]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	var request authRequest.ValidateTokenRequest

	// validate incoming request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// calling gRPC sso-service ValidateToken handler
	response, err := h.SSOClient.ValidateToken(request.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authResponse.ValidateTokenResponse{
		Valid: response.Valid,
		UserId: response.UserId,
	})
}


func (h *AuthHandler) GetUserByToken(c *gin.Context) {
	var request struct {
		Token string `json:"token" binding:"required"`
	}

	// validate incoming request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// calling gRPC sso-service GetUserByToken handler
	response, err := h.SSOClient.GetUserByToken(request.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": response.UserId,
		"login": response.Login,
		"username": response.Username,
	})
}

// @Summary Setup Google 2FA
// @Description Setup 2Fa using Google Authenticator
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} authResponse.Setup2FAResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /2fa/setup [post]
func (h *AuthHandler) Setup2FA(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	qrURL, err := h.SSOClient.Setup2FA(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, authResponse.Setup2FAResponse{
		QrURL: qrURL,
	})
}

// @Summary Verify 2FA code
// @Description Verify 2FA code to enable 2FA for account
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body authRequest.Verify2FARequest true "2FA code"
// @Success 200 {object} authResponse.Verify2FAResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /2fa/verify [post]
func (h *AuthHandler) Verify2FA(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "userID missed"})
		return
	}
	var request authRequest.Verify2FARequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	verif, err := h.SSOClient.Verify2FA(userID, request.Code)
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}
	if !verif {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "wrong 2FA code"})
		return
	}
	c.JSON(http.StatusOK, authResponse.Verify2FAResponse{})
}