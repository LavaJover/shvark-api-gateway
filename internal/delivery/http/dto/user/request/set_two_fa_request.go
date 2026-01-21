package request

type SetTwoFaSecretRequest struct {
	TwoFaSecret string `json:"two_fa_secret" binding:"required"`
}

type SetTwoFaEnabledRequest struct {
	Enabled bool `json:"enabled" binding:"required"`
}