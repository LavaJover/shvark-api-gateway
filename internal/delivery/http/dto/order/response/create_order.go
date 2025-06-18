package response

import "time"

type CreateOrderResponse struct {
	OrderID			string 			`json:"order_id"`
	OrderStatus 	string 			`json:"order_status"`
	BankDetail  	BankDetail		`json:"bank_detail"`
	AmountFiat 		float64 		`json:"amount_fiat"`
	AmountCrypto    float64  		`json:"amount_crypto"`
}

type BankDetail struct {
	ID 		 		string 	 `json:"id"`
	TraderID 		string 	 `json:"trader_id"`
	Currency 		string 	 `json:"currency"`
	Country  		string 	 `json:"country"`
	MinAmount 		float64  `json:"min_amount"`
	MaxAmount 		float64  `json:"max_amount"`
	BankName 		string 	 `json:"bank_name"`
	PaymentSystem 	string   `json:"payment_system"`
	Enabled 		bool 	 `json:"enabled"`
	Delay 			string 	 `json:"delay"`
	Owner			string   `json:"owner"`
	CardNumber      string   `json:"card_number"`
	Phone 			string   `json:"phone"`
}

type Order struct {
	OrderID 	 string     `json:"order_id"`
	Status  	 string     `json:"status"`
	AmountFiat 	 float64    `json:"amount_fiat"`
	AmountCrypto float64    `json:"amount_crypto"`
	ExpiresAt    time.Time  `json:"expires_at"`
	BankDetail   BankDetail `json:"bank_detail"`
	TraderReward float64	`json:"trader_reward"`
	CreatedAt 	 time.Time	`json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}