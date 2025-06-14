package request

type CreateH2HPayInRequest struct {
	MerchantID      string  `json:"merchant_id"`
	Currency 		string 	`json:"currency"`
	PaymentSystem 	string 	`json:"payment_system"`
	AmountFiat   	float64 `json:"amount_fiat"`
	MerchantOrderID string  `json:"merchant_order_id"`
	FTD				bool    `json:"ftd"`
	CallbackURL 	string  `json:"callback_url"`
	ClientID		string  `json:"client_id"`
	Shuffle			int32   `json:"shuffle"`
	TTL 			string	`json:"ttl"`
}