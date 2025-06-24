package response

type LoginResponse struct{
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TwoFaRequiredResponse struct {
	TwoFaRequired bool `json:"two_fa_required"`
}