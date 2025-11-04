package request

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	TwoFaCode string `json:"twoFaCode"`
}
