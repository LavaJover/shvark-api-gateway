package response

import "time"

type GetTraderHistoryResponse struct {
	History    []Transaction `json:"history"`
	Pagination Pagination    `json:"pagination"`
}

type Pagination struct {
	TotalItems   int  `json:"totalItems"`
	TotalPages   int  `json:"totalPages"`
	CurrentPage  int  `json:"currentPage"`
	ItemsPerPage int  `json:"itemsPerPage"`
	HasNextPage  bool `json:"hasNextPage"`
	HasPrevPage  bool `json:"hasPrevPage"`
}

type Transaction struct {
	ID       string  `json:"id"`
	TraderID string  `json:"traderId"`
	Currency string  `json:"currency"`
	Type     string  `json:"type"`
	Amount   float64 `json:"amount"`
	OrderID  string  `json:"orderId"`
	TxHash   string  `json:"txHash"`
	Status   string  `json:"status"`
	CreatedAt time.Time `json:"createdAt"` // Добавьте это поле
}

type GetTraderHistoryErrorResponse struct {
	Error string `json:"error"`
}