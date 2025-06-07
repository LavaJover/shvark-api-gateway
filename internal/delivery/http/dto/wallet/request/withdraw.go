package request

type WithdrawRequest struct {
	TraderID	string `json:"traderId"`
	ToAddress 	string `json:"toAddress"`
	Amount float64	`json:"amount"`
}