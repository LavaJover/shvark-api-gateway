package handlers

import (
	"net/http"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	UserClient *client.UserClient
}

func NewUserHandler(addr string) (*UserHandler, error) {
	userClient, err := client.NewUserClient(addr)
	if err != nil {
		return nil, err
	}
	return &UserHandler{
		UserClient: userClient,
	}, nil
}

// @Summary Get user by UUID
// @Description Get user by UUID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} GetUserByIDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.UserClient.GetUserByID(userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": response.UserId,
		"login": response.Login,
		"username": response.Username,
		"password": response.Password,
	})
}

type GetUserByIDResponse struct {
	UserID string `json:"user_id"`
	Login string `json:"login"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	
}

func (h *UserHandler) UpdateUser(c *gin.Context) {

}

func (h *UserHandler) DeleteUser(c *gin.Context) {

}