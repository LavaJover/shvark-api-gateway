package request

type WithdrawRequest struct {
	Amount 	 	float64 `json:"amount"`
	Currency 	string 	`json:"currency"`
	ToAddress 	string  `json:"toAddress"`
}

type WithdrawResponse struct {
	TxHash 		string 	`json:"txHash"`
	Amount 		float64 `json:"amount"`
	Currency 	string 	`json:"currency"`
	ToAddress 	string 	`json:"toAddress"`
}