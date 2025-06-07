package response

type GetTraderBalanceResponse struct {
	TraderID string `json:"traderId"`
	Currency string `json:"currency"`
	Balance float64 `json:"balance"`
	Frozen 	float64 `json:"frozen"`
	Address string 	`json:"address"`
}

type GetTraderBalanceErrorResponse struct {
	Error string `json:"error"`
}