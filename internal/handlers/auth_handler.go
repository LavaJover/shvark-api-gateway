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
// @Param input body RegisterRequest true "Registration credentials"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var request struct {
		Login 	 string	`json:"login" binding:"required"`
		Username string	`json:"username" binding:"required"`
		Password string	`json:"password" binding:"required,min=8"`
	}

	// validate incoming request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// calling gRPC sso-service
	resp, err := h.SSOClient.Register(request.Login, request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": resp.UserId,
		"message": resp.Message,
	})
}

type RegisterRequest struct {
	Login string `json:"login" example:"CoolUserLogin"`
	Username string `json:"username" example:"CoolUsername"`
	Password string `json:"password" example:"securepass123"`
}

type RegisterResponse struct {
	UserID string `json:"user_id" example:"1d6ab366-4fca-4bcc-972d-875c35ea939a"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid email"`
}

// @Summary User login
// @Description Log in user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var request struct{
		Login 		string `json:"login" binding:"required"`
		Password 	string `json:"password" binding:"required,min=8"`
	}

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

	c.JSON(http.StatusOK, gin.H{
		"access_token": response.AccessToken,
		"refresh_token": response.RefreshToken,
	})
}

type LoginRequest struct {
	Login 	 string `json:"login" example:"CoolUserLogin"`
	Password string `json:"password" example:"securepass123"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary JWT validation check-point
// @Description Check if JWT is valid or not
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body ValidateTokenRequest true "JWT instance"
// @Success 200 {object} ValidateTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/validate_token [post]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	var request struct {
		Token string `json:"token" binding:"required"`
	}

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

	c.JSON(http.StatusOK, gin.H{
		"valid": response.Valid,
		"user_id": response.UserId,
	})
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Valid  bool 	`json:"valid"`
	UserId string 	`json:"user_id" example:"1d6ab366-4fca-4bcc-972d-875c35ea939a"`
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