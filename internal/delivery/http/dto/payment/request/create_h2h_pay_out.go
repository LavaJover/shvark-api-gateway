package request

type CreateH2HPayOutRequest struct {
	Currency string `json:"currency"`
	PaymentSystem string `json:"payment_system"`
	BankName string `json:"bank_name"`
	Amount float64 `json:"amount"`
	MerchantOrderID string `json:"merchant_order_id"`
	Ftd bool `json:"ftd"`
	PaymentDetails PaymentDetails `json:"payment_details"`
	CallbackURL string `json:"callback_url"`
	MerchantID string `json:"merchant_id"`
}

type PaymentDetails struct {
	Phone string `json:"phone"`
	Bank string `json:"bank"`
	CardNumber string `json:"card_number"`
}