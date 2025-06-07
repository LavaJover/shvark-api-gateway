package response

type GetTraderHistoryResponse struct {
	History []Transaction `json:"history"`
}

type Transaction struct {
	TraderID string `json:"traderId"`
	Currency string `json:"currency"`
	Type 	 string `json:"type"`
	Amount 	 float64 `json:"amount"`
	OrderID  string  `json:"orderId"`
	Status 	 string  `json:"status"`
}