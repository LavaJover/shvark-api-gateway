package response

type LoginResponse struct{
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	DateTimeExpires string `json:"dateTimeExpires"`
}

type TwoFaRequiredResponse struct {
	TwoFaRequired bool `json:"two_fa_required"`
}