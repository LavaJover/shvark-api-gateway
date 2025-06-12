package request

type CancelOrderRequest struct {
	OrderID string `json:"order_id"`
}