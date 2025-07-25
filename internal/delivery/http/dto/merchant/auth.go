package merchant

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	TwoFaCode string `json:"twoFaCode"`
}

type LoginResponse struct {
	Token string `json:"token"`
	DateTimeExpires string `json:"dateTimeExpires"`
	ApiKey string `json:"apiKey"`
}