package response

type CreateH2HPayOutResponse struct {
	ID string `json:"id"`
	Currency string `json:"currency"`
	Amount float64 `json:"amount"`
	UsdAmount float64 `json:"usd_amount"`
	UsdRate float64 `json:"usd_rate"`
	PaymentSystem string `json:"payment_system"`
	Status string `json:"status"`
	MerchantOrderID string `json:"merchant_order_id"`
	CallbackURL string `json:"callback_url"`
	ExpiresAt int64 `json:"expires_at"`
}