package request

type CreateH2HPayInRequest struct {
	StoreID			string  `json:"store_id"`
	Currency 		string 	`json:"currency"`
	PaymentSystem 	string 	`json:"paymentSystem"`
	AmountFiat   	float64 `json:"amountFiat"`
	MerchantOrderID string  `json:"merchantOrderId"`
	FTD				bool    `json:"ftd"`
	CallbackURL 	string  `json:"callbackUrl"`
	ClientID		string  `json:"clientId"`
	Shuffle			int32   `json:"shuffle"`
	Issuer			string  `json:"issuer"`
}