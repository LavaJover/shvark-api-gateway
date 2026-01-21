package handlers

import (
	"net/http"

	userservice "github.com/LavaJover/shvark-api-gateway/internal/client/user-service"
	userRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/user/request"
	userResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/user/response"
	"github.com/LavaJover/shvark-api-gateway/internal/domain"
	userpb "github.com/LavaJover/shvark-user-service/proto/gen"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	UserClient *userservice.UserClient
}

func NewUserHandler(addr string) (*UserHandler, error) {
	userClient, err := userservice.NewUserClient(addr)
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
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} userResponse.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	response, err := h.UserClient.GetUserByID(userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// GetUserByIDResponse не содержит поле Role, поэтому мы не можем его использовать
	// Если нужно получить роль, нужно использовать другой метод или обновить прото
	c.JSON(http.StatusOK, userResponse.UserResponse{
		ID:           response.UserId,
		Login:        response.Login,
		Username:     response.Username,
		TwoFaEnabled: response.TwoFaEnabled,
		// Role не доступен в GetUserByIDResponse
	})
}

// @Summary Create new user
// @Description Create a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userRequest.CreateUserRequest true "User data"
// @Success 201 {object} userResponse.CreateUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req userRequest.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.UserClient.CreateUser(req.Login, req.Username, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		// Проверяем тип ошибки (в реальном коде нужно использовать константы)
		if err.Error() == "user already exists" || err.Error() == "already exists" {
			status = http.StatusConflict
		} else if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, userResponse.CreateUserResponse{
		UserID: response.UserId,
	})
}

// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body userRequest.UpdateUserRequest true "User data to update"
// @Success 200 {object} userResponse.UpdateUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req userRequest.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создаем объект для обновления
	updateData := &domain.UpdateUserData{}
	
	// Собираем поля для обновления
	updateFields := []string{}

	if req.Login != nil {
		updateData.Login = *req.Login
		updateFields = append(updateFields, "login")
	}
	if req.Username != nil {
		updateData.Username = *req.Username
		updateFields = append(updateFields, "username")
	}
	if req.Password != nil {
		updateData.Password = *req.Password
		updateFields = append(updateFields, "password")
	}
	if req.Role != nil {
		updateData.Role = *req.Role
		updateFields = append(updateFields, "role")
	}
	if req.TwoFaSecret != nil {
		updateData.TwoFaSecret = *req.TwoFaSecret
		updateFields = append(updateFields, "two_fa_secret")
	}
	if req.TwoFaEnabled != nil {
		updateData.TwoFaEnabled = *req.TwoFaEnabled
		updateFields = append(updateFields, "two_fa_enabled")
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	response, err := h.UserClient.UpdateUser(userID.String(), updateData, updateFields)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userResponse.UpdateUserResponse{
		User: userResponse.UserResponse{
			ID:           response.User.UserId,
			Login:        response.User.Login,
			Username:     response.User.Username,
			Role:         response.User.Role,
			TwoFaEnabled: response.User.TwoFaEnabled,
		},
	})
}

// @Summary Delete user
// @Description Delete user from the system
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} userResponse.DeleteUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	response, err := h.UserClient.DeleteUser(&userpb.DeleteUserRequest{
		UserId: userID.String(),
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userResponse.DeleteUserResponse{
		Success: response.Success,
	})
}

// @Summary Get users with filters
// @Description Get users with filtering, sorting and pagination
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" minimum(1) default(1)
// @Param limit query int false "Items per page" minimum(1) maximum(100) default(20)
// @Param login_filter query string false "Filter by login (substring)"
// @Param username_filter query string false "Filter by username (substring)"
// @Param roles query []string false "Filter by roles" collectionFormat(multi) Enums(TRADER, MERCHANT, ADMIN, TEAM_LEAD)
// @Param sort_by query string false "Sort field" Enums(login, username, created_at, role)
// @Param sort_desc query boolean false "Sort descending" default(false)
// @Success 200 {object} userResponse.GetUsersResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetUsersWithFilter(c *gin.Context) {
	var req userRequest.FilterUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем значения по умолчанию
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	response, err := h.UserClient.GetUsersWithFilter(&userpb.GetUsersWithFilterRequest{
		Page:           req.Page,
		Limit:          req.Limit,
		LoginFilter:    req.LoginFilter,
		UsernameFilter: req.UsernameFilter,
		Roles:          req.Roles,
		SortBy:         req.SortBy,
		SortDesc:       req.SortDesc,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	users := make([]userResponse.UserResponse, len(response.Users))
	for i, user := range response.Users {
		users[i] = userResponse.UserResponse{
			ID:           user.UserId,
			Login:        user.Login,
			Username:     user.Username,
			Role:         user.Role,
			TwoFaEnabled: user.TwoFaEnabled,
		}
	}

	// Приводим int32 к int64
	c.JSON(http.StatusOK, userResponse.GetUsersResponse{
		TotalPages: int64(response.TotalPages),
		TotalCount: response.TotalCount,
		Users:      users,
	})
}

// @Summary Change user password
// @Description Change user's password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body userRequest.ChangePasswordRequest true "New password"
// @Success 200 {object} userResponse.ChangePasswordResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req userRequest.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.UserClient.ChangePassword(&userpb.ChangePasswordRequest{
		UserId:      userID.String(),
		NewPassword: req.NewPassword,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userResponse.ChangePasswordResponse{
		Success: response.Success,
	})
}

type SuccessResponse struct {
	Message string
}

// @Summary Set 2FA secret
// @Description Set 2FA secret for user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body userRequest.SetTwoFaSecretRequest true "2FA secret"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/2fa/secret [post]
func (h *UserHandler) SetTwoFaSecret(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req userRequest.SetTwoFaSecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.UserClient.SetTwoFaSecret(&userpb.SetTwoFaSecretRequest{
		UserId:      userID.String(),
		TwoFaSecret: req.TwoFaSecret,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// @Summary Get 2FA secret
// @Description Get 2FA secret for user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} userResponse.TwoFaSecretResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/2fa/secret [get]
func (h *UserHandler) GetTwoFaSecret(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	response, err := h.UserClient.GetTwoFaSecretByID(&userpb.GetTwoFaSecretByIDRequest{
		UserId: userID.String(),
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userResponse.TwoFaSecretResponse{
		TwoFaSecret: response.TwoFaSecret,
	})
}

// @Summary Enable/disable 2FA
// @Description Enable or disable 2FA for user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body userRequest.SetTwoFaEnabledRequest true "Enable/disable flag"
// @Success 200 {object} userResponse.SetTwoFaEnabledResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/2fa [put]
func (h *UserHandler) SetTwoFaEnabled(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req userRequest.SetTwoFaEnabledRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Не присваиваем результат переменной, если он не используется
	_, err = h.UserClient.SetTwoFaEnabled(&userpb.SetTwoFaEnabledRequest{
		UserId:  userID.String(),
		Enabled: req.Enabled,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userResponse.SetTwoFaEnabledResponse{
		Success: true,
	})
}

// @Summary Get traders by team lead
// @Description Get all traders in a team lead's team
// @Tags teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param team_lead_id path string true "Team Lead ID"
// @Success 200 {object} []userResponse.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams/{team_lead_id}/traders [get]
func (h *UserHandler) GetTeamTraders(c *gin.Context) {
	teamLeadID, err := uuid.Parse(c.Param("team_lead_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team lead ID format"})
		return
	}

	response, err := h.UserClient.GetTeamTraders(&userpb.GetTeamTradersRequest{
		TeamLeadId: teamLeadID.String(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	traders := make([]userResponse.UserResponse, len(response.Traders))
	for i, trader := range response.Traders {
		traders[i] = userResponse.UserResponse{
			ID:       trader.UserId,
			Login:    trader.Login,
			Username: trader.Username,
			Role:     trader.Role,
		}
	}

	c.JSON(http.StatusOK, traders)
}

// @Summary Add trader to team
// @Description Add a trader to a team lead's team
// @Tags teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param team_lead_id path string true "Team Lead ID"
// @Param request body userRequest.AddTraderToTeamRequest true "Trader data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams/{team_lead_id}/traders [post]
func (h *UserHandler) AddTraderToTeam(c *gin.Context) {
	teamLeadID, err := uuid.Parse(c.Param("team_lead_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team lead ID format"})
		return
	}

	var req userRequest.AddTraderToTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.UserClient.AddTraderToTeam(&userpb.AddTraderToTeamRequest{
		TeamLeadId: teamLeadID.String(),
		TraderId:   req.TraderID,
		Commission: req.Commission,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// @Summary Update relationship parameters
// @Description Update commission parameters for a team relationship
// @Tags teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param relationship_id path string true "Relationship ID"
// @Param request body userRequest.UpdateRelationshipParamsRequest true "New parameters"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /relationships/{relationship_id} [put]
func (h *UserHandler) UpdateRelationshipParams(c *gin.Context) {
	relationshipID, err := uuid.Parse(c.Param("relationship_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid relationship ID format"})
		return
	}

	var req userRequest.UpdateRelationshipParamsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.UserClient.UpdateRelationshipParams(&userpb.UpdateRelationshipParamsRequest{
		Relationship: &userpb.Relationship{
			Id:         relationshipID.String(),
			Commission: req.Commission,
		},
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "relationship not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// @Summary Promote user to team lead
// @Description Promote a user to team lead role
// @Tags teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/promote-to-team-lead [post]
func (h *UserHandler) PromoteToTeamLead(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	_, err = h.UserClient.PromoteToTeamLead(&userpb.PromoteToTeamLeadRequest{
		UserId: userID.String(),
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// @Summary Demote team lead
// @Description Demote a team lead to trader role
// @Tags teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Team Lead ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/demote-team-lead [post]
func (h *UserHandler) DemoteTeamLead(c *gin.Context) {
	teamLeadID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team lead ID format"})
		return
	}

	_, err = h.UserClient.DemoteTeamLead(&userpb.DemoteTeamLeadRequest{
		TeamLeadId: teamLeadID.String(),
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}