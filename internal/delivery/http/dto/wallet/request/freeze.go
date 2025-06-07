package request

type FreezeRequest struct {
	TraderID string `json:"traderId"`
	Amount	float64 `json:"amount"`
	OrderID string `json:"orderId"`
}