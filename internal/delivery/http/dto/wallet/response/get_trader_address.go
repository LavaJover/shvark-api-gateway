package response

type GetTraderWalletAddressResponse struct {
	Address string `json:"address"`
}

type GetTraderWalletAddressErrorResponse struct {
	Error string `json:"error"`
}