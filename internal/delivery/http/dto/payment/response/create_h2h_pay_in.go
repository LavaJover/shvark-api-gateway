package response

type CreateH2HPayInResponse struct {
	OrderID 		string  	   `json:"order_id"`
	AmountFiat 		float64		   `json:"amount_fiat"`
	AmountCrypto 	float64		   `json:"amount_crypto"`
	Currency 		string 		   `json:"currency"`
	PaymentSystem	string		   `json:"payment_system"`
	Status 			string 		   `json:"status"`
	MerchantOrderID	string		   `json:"merchant_order_id"`
	CallbackURL		string		   `json:"callback_url"`
	PaymentDetails 	PaymentDetails `json:"payment_details"`
	ExpiresAt 		int64 		   `json:"expires_at"`
	TPayLink 		string 		   `json:"tpay_link"`
}

type PaymentDetails struct {
	CardNumber string `json:"card_number"`
	Owner	   string `json:"owner"`
	Phone 	   string `json:"phone"`
	BankID     string `json:"bank_id"` // In SBP system
	Bank 	   string `json:"bank"` // User-friendly
	BankName   string `json:"bank_name"`
}

type NoBankDetailsErrorResponse struct {
	Error string `json:"error"`
}

type BadRequestErrorResponse struct {
	Error string `json:"error"`
}