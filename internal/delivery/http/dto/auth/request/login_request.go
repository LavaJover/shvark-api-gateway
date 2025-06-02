package request

type LoginRequest struct {
	Login 	 string `json:"login" example:"CoolUserLogin"`
	Password string `json:"password" example:"securepass123"`
}