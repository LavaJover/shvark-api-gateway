package response

type AssignRoleResponse struct {
	Success bool `json:"success" binding:"required"`
}