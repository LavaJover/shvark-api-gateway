
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
	"github.com/gin-gonic/gin"
	authRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/auth/request"
	authResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/auth/response"
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
// @Tags Auth
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
	response, err := h.SSOClient.Register(request.Login, request.Username, request.Password)
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
// @Tags Auth
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
	response, err := h.SSOClient.Login(request.Login, request.Password)
	if err != nil {
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
// @Tags Auth
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