package merchant

type GetAccountBalanceRequest struct {}

type GetAccountBalanceResponse struct {
	Balances []AccountBalance `json:"balances"`
}

type AccountBalance struct {
	Name 		string 	`json:"name"`
	Balance 	string 	`json:"balance"`
	Currency 	string  `json:"currency"`
}