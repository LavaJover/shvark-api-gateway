package request

type OffchainWithdrawRequest struct {
	TraderID string `json:"traderId"`
	Amount   float64 `json:"amount"`
	TxHash	 string  `json:"txHash"`
}