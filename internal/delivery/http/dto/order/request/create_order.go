package request

type CreateOrderRequest struct {
	MerchantID	  string 		`json:"merchant_id"`
	AmountFiat    float64 		`json:"amount_fiat"`
	Currency	  string		`json:"currency"`
	Country		  string		`json:"country"`
	ClientData	  string 		`json:"client_data"`
	PaymentSystem string 		`json:"payment_system"`
	Metadata      string 		`json:"metadata"`
	TTL 		  string 		`json:"ttl"`
}