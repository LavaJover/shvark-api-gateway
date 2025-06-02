package response

type CreateOrderResponse struct {
	OrderID		string `json:"order_id"`
	OrderStatus string `json:"order_status"`
	BankDetail  BankDetail
	Amount 		float64 `json:"amount"`
}

type BankDetail struct {
	ID 		 string 	 `json:"id"`
	TraderID string 	 `json:"trader_id"`
	Currency string 	 `json:"currency"`
	Country  string 	 `json:"country"`
	MinAmount float64 	 `json:"min_amount"`
	MaxAmount float64 	 `json:"max_amount"`
	BankName string 	 `json:"bank_name"`
	PaymentSystem string `json:"payment_system"`
	Enabled bool 		 `json:"enabled"`
	Delay string 		 `json:"delay"`
}