package response

type WithdrawResponse struct {
	TxID string `json:"txid"`
}

type WithdrawErrorResponse struct {
	Error 	string `json:"error"`
	Details string `json:"details"`
}