package response

type GetTraderHistoryResponse struct {
	History []Transaction `json:"history"`
}

type Transaction struct {
	ID		 string `json:"id"`
	TraderID string `json:"traderId"`
	Currency string `json:"currency"`
	Type 	 string `json:"type"`
	Amount 	 float64 `json:"amount"`
	OrderID  string  `json:"orderId"`
	TxHash   string  `json:"txHash"`
	Status 	 string  `json:"status"`
	Metadata string  `json:"metadata"`
}

type GetTraderHistoryErrorResponse struct {
	Error string `json:"error"`
}