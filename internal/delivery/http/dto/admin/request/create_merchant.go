package request

type CreateMerchantRequest struct {
	Username string `json:"username"`
	Login string `json:"login"`
	Password string `json:"password"`
}