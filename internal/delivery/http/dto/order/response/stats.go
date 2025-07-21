package response

type GetOrderStatsResponse struct {
	TotalOrders int64 `json:"total_orders"`
	SucceedOrders int64 `json:"succeed_orders"`
	CanceledOrders int64 `json:"canceled_orders"`
	ProcessedAmountFiat float64 `json:"processed_amount_fiat"`
	ProcessedAmountCrypto float64 `json:"processed_amount_crypto"`
	CanceledAmountFiat float64 `json:"canceled_amount_fiat"`
	CanceledAmountCrypto float64 `json:"canceled_amount_crypto"`
	IncomeCrypto float64 `json:"income_crypto"`
}