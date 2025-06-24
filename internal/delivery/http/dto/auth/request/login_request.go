package request

type LoginRequest struct {
	Login 	 string 	`json:"login" example:"CoolUserLogin"`
	Password string 	`json:"password" example:"securepass123"`
	TwoFACode string 	`json:"two_fa_code"`
}