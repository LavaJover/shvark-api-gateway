package handlers

import (
	"net/http"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	"github.com/gin-gonic/gin"
	authzRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/authz/request"
	authzResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/authz/response"
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
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body authzRequest.AssignRoleRequest true "Role assigned to user"
// @Success 200 {object} authzResponse.AssignRoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rbac/roles [post]
func (h *AuthzHandler) AssignRole(c *gin.Context) {
	var request authzRequest.AssignRoleRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.AuthzClient.AssignRole(request.UserID, request.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authzResponse.AssignRoleResponse{
		Success: response.Success,
	})
}

// @Summary Revoke role
// @Description Revoke role
// @Tags RBAC
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body authzRequest.RevokeRoleRequest true "Role assigned to user to revoke"
// @Success 200 {object} authzResponse.RevokeRoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rbac/roles [delete]
func (h *AuthzHandler) RevokeRole(c *gin.Context) {
	var request authzRequest.AssignRoleRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.AuthzClient.RevokeRole(request.UserID, request.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authzResponse.RevokeRoleResponse{
		Success: response.Success,
	})
}

// @Summary Add policy
// @Description Add policy
// @Tags RBAC
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body authzRequest.AddPolicyRequest true "New policy details"
// @Success 200 {object} authzResponse.AddPolicyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rbac/policies [post]
func (h *AuthzHandler) AddPolicy(c *gin.Context) {
	var request authzRequest.AddPolicyRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.AuthzClient.AddPolicy(request.Role, request.Object, request.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authzResponse.AddPolicyResponse{
		Success: response.Success,
	})
}

// @Summary Delete policy
// @Description Delete policy
// @Tags RBAC
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body authzRequest.DeletePolicyRequest true "Policy details to delete"
// @Success 200 {object} authzResponse.DeletePolicyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rbac/policies [delete]
func (h *AuthzHandler) DeletePolicy(c *gin.Context) {
	var request authzRequest.DeletePolicyRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.AuthzClient.DeletePolicy(request.Role, request.Object, request.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authzResponse.DeletePolicyResponse{
		Success: response.Success,
	})
}

// @Summary Check user permission
// @Description Check user permission
// @Tags RBAC
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body authzRequest.CheckPermissionRequest true "Permission subject, object, action"
// @Success 200 {object} authzResponse.CheckPermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rbac/permissions [post]
func (h *AuthzHandler) CheckPermission(c *gin.Context) {
	var request authzRequest.CheckPermissionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.AuthzClient.CheckPermission(request.UserID, request.Object, request.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK,authzResponse.CheckPermissionResponse{
		Allowed: response.Allowed,
	})
}