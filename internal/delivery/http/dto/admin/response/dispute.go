package response

import "time"

type CreateDisputeResponse struct {
	DisputeID string `json:"dispute_id"`
}

type AcceptDisputeResponse struct {

}

type RejectDisputeResponse struct {

}

type GetDisputeInfoResponse struct {
	Dispute Dispute `json:"dispute"`
}

type FreezeDisputeResponse struct {

}

type BankDetail struct {
	BankName 		string `json:"bank_name"`
	PaymentSystem 	string `json:"payment_system"`
	Phone 			string `json:"phone"`
	CardNumber 		string `json:"card_number"`
	Owner 			string `json:"owner"`
	TraderID 		string `json:"trader_id"`
}

type Order struct {
	ID 				string 		`json:"order_id"`
	MerchantOrderID string 		`json:"merchant_order_id"`
	AmountFiat 		float64 	`json:"amount_fiat"`
	CryproRate 		float64 	`json:"crypro_rate"`
	AmountCrypto 	float64		`json:"amount_crypto"`
	BankDetail 		BankDetail 	`json:"bank_detail"`
}

type Dispute struct {
	DisputeID 	  		string 		`json:"dispute_id"`
	ProofUrl 	  		string 		`json:"proof_url"`
	DisputeReason 		string 		`json:"dispute_reason"`
	DisputeStatus 		string 		`json:"dispute_status"`
	DisputeAmountFiat 	float64 	`json:"dispute_amount_fiat"`
	DisputeAmountCrypto float64 	`json:"dispute_amount_crypto"`
	DisputeCryptoRate 	float64		`json:"dispute_crypto_rate"`
	Order 				Order 		`json:"order"`
	OrderID 			string 		`json:"order_id"`
	AcceptAt			time.Time 	`json:"accept_at"`
}

type GetOrderDisputesResponse struct {
	Disputes []Dispute `json:"disputes"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	CurrentPage  int64 `json:"current_page"`
	TotalPages 	 int64 `json:"total_pages"`
	TotalItems 	 int64 `json:"total_items"`
	ItemsPerPage int64 `json:"items_per_page"`
}