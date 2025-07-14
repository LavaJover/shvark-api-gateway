package request

type RegisterRequest struct {
	Login 	 string	`json:"login" binding:"required"`
	Username string	`json:"username" binding:"required"`
	Password string	`json:"password" binding:"required,min=8"`
	Role 	 string `json:"role"`
}