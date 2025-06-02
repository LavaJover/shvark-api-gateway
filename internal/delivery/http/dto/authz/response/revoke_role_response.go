package response

type RevokeRoleResponse struct {
	Success bool `json:"success" binding:"required"`
}