package response

type LoginResponse struct {
	Token string `json:"token"`
	DateTimeExpires string `json:"dateTimeExpires"`
	ApiKey string `json:"apiKey"`
}