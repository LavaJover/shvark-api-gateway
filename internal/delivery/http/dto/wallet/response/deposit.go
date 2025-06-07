package response 

type DepositResponse struct {
	Success bool `json:"success"`
}

type DepositErrorResponse struct {
	Error string `json:"error" example:"Wallet not found"`
}