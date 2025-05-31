package handlers

import (
	"net/http"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	"github.com/gin-gonic/gin"
)

type AuthzHandler struct {
	AuthzClient *client.AuthzClient
}

func NewAuthzhandler(addr string) (*AuthzHandler, error) {
	authzClient, err := client.NewAuthzClient(addr)
	if err != nil {
		return nil, err
	}
	return &AuthzHandler{
		AuthzClient: authzClient,
	}, nil
}

// @Summary Assign role
// @Description Assign role
// @Tags RBAC
// @Accept json
// @Produce json
// @Param input body AssignRoleRequest true "Role assigned to user"
// @Success 200 {object} AssignRoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rbac/roles [post]
func (h *AuthzHandler) AssignRole(c *gin.Context) {
	var request AssignRoleRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.AuthzClient.AssignRole(request.UserID, request.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": response.Success,
	})
}

type AssignRoleRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role string `json:"role" binding:"required"`
}

type AssignRoleResponse struct {
	Success bool `json:"success" binding:"required"`
}

// @Summary Revoke role
// @Description Revoke role
// @Tags RBAC
// @Accept json
// @Produce json
// @Param input body RevokeRoleRequest true "Role assigned to user to revoke"
// @Success 200 {object} RevokeRoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rbac/roles [delete]
func (h *AuthzHandler) RevokeRole(c *gin.Context) {
	var request RevokeRoleRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.AuthzClient.RevokeRole(request.UserID, request.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": response.Success,
	})
}

type RevokeRoleRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role string `json:"role" binding:"required"`
}

type RevokeRoleResponse struct {
	Success bool `json:"success" binding:"required"`
}

// @Summary Add policy
// @Description Add policy
// @Tags RBAC
// @Accept json
// @Produce json
// @Param input body AddPolicyRequest true "New policy details"
// @Success 200 {object} AddPolicyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rbac/policies [post]
func (h *AuthzHandler) AddPolicy(c *gin.Context) {
	var request AddPolicyRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.AuthzClient.AddPolicy(request.Role, request.Object, request.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": response.Success,
	})
}

type AddPolicyRequest struct {
	Role string `json:"role" binding:"required"`
	Object string `json:"object" binding:"required"`
	Action string `json:"action" binding:"required"`
}

type AddPolicyResponse struct {
	Success bool `json:"success" binding:"required"`
}

// @Summary Delete policy
// @Description Delete policy
// @Tags RBAC
// @Accept json
// @Produce json
// @Param input body DeletePolicyRequest true "Policy details to delete"
// @Success 200 {object} DeletePolicyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rbac/policies [delete]
func (h *AuthzHandler) DeletePolicy(c *gin.Context) {
	var request DeletePolicyRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.AuthzClient.DeletePolicy(request.Role, request.Object, request.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": response.Success,
	})
}

type DeletePolicyRequest struct {
	Role string `json:"role" binding:"required"`
	Object string `json:"object" binding:"required"`
	Action string `json:"action" binding:"required"`
}

type DeletePolicyResponse struct {
	Success bool `json:"success" binding:"required"`
}

// @Summary Check user permission
// @Description Check user permission
// @Tags RBAC
// @Accept json
// @Produce json
// @Param input body CheckPermissionRequest true "Permission subject, object, action"
// @Success 200 {object} CheckPermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rbac/permissions [post]
func (h *AuthzHandler) CheckPermission(c *gin.Context) {
	var request CheckPermissionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.AuthzClient.CheckPermission(request.UserID, request.Object, request.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"allowed": response.Allowed,
	})
}

type CheckPermissionRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Object string `json:"object" binding:"required"`
	Action string `json:"action" binding:"required"`
}

type CheckPermissionResponse struct {
	Allowed bool `json:"allowed" binding:"required"`
}