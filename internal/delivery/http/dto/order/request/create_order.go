package request

type CreateOrderRequest struct {
	MerchantID	  string 	`json:"merchant_id"`
	Amount		  float64 	`json:"amount"`
	Currency	  string	`json:"currency"`
	Country		  string	`json:"country"`
	ClientData	  string 	`json:"client_data"`
	PaymentSystem string 	`json:"payment_system"`
	Metadata      string 	`json:"metadata"`
}