package request

type CreateDeeplinkRequest struct {
	OrderID     string  `json:"order_id" binding:"required"`
	BankCode    string  `json:"bank_code" binding:"required"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}