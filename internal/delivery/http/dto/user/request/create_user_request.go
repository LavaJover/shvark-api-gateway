package request

type CreateUserRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=50"`
	Username string `json:"username" binding:"required,min=2,max=100"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Role     string `json:"role" binding:"required,oneof=TRADER MERCHANT ADMIN"`
}