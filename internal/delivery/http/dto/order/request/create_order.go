package request

type CreateOrderRequest struct {
	MerchantOrderID string      `json:"merchant_order_id"`
	MerchantID		string		`json:"merchant_id"`
	AmountFiat    	float64 	`json:"amount_fiat"`
	Currency	  	string		`json:"currency"`
	Country		  	string		`json:"country"`
	ClientID	  	string 		`json:"client_id"`
	PaymentSystem 	string 		`json:"payment_system"`
	TTL 		  	string 		`json:"ttl"`
	Shuffle 		int32		`json:"shuffle"`
}