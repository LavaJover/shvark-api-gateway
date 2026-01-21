package request

type ChangePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6,max=100"`
}