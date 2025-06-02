package request

type AssignRoleRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role string `json:"role" binding:"required"`
}