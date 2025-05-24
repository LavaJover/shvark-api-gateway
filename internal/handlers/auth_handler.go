package handlers

import (
	"net/http"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	ssoClient *client.SSOClient
}

func NewAuthHandler(addr string) (*AuthHandler, error) {
	ssoClient, err := client.NewSSOClient(addr)
	if err != nil {
		return nil, err
	}

	return &AuthHandler{ssoClient: ssoClient}, nil
}

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
	resp, err := h.ssoClient.Register(request.Login, request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": resp.UserId,
		"message": resp.Message,
	})
}

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
	response, err := h.ssoClient.Login(request.Login, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": response.AccessToken,
		"refresh_token": response.RefreshToken,
	})
}

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
	response, err := h.ssoClient.ValidateToken(request.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": response.Valid,
		"user_id": response.UserId,
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
	response, err := h.ssoClient.GetUserByToken(request.Token)
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