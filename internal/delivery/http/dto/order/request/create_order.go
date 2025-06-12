package request

import "time"

type CreateOrderRequest struct {
	MerchantID	  string 	`json:"merchant_id"`
	AmountFiat    float64 	`json:"amount_fiat"`
	Currency	  string	`json:"currency"`
	Country		  string	`json:"country"`
	ClientData	  string 	`json:"client_data"`
	PaymentSystem string 	`json:"payment_system"`
	Metadata      string 	`json:"metadata"`
	ExpiresAt     time.Time `json:"expires_at"`
}