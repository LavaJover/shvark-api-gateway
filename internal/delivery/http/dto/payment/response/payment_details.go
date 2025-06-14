package response

type PaymentDetails struct {
	CardNumber string `json:"card_number"`
	Owner	   string `json:"owner"`
	Phone 	   string `json:"phone"`
	BankID     string `json:"bank_id"` // In SBP system
	Bank 	   string `json:"bank"` // User-friendly
	BankName   string `json:"bank_name"`
}