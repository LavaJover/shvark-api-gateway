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
		Password string	`json:"password" binding:"required, min=8"`
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