package response

type CheckPermissionResponse struct {
	Allowed bool `json:"allowed" binding:"required"`
}