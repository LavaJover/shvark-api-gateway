package request

type CheckPermissionRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Object string `json:"object" binding:"required"`
	Action string `json:"action" binding:"required"`
}