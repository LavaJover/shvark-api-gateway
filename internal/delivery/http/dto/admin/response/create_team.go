package response

type CreateTeamResponse struct {
	TraderID 		string 		`json:"trader_id"` 
	AccessToken 	string 		`json:"access_token"`
	WalletAddress 	string 		`json:"wallet_address"`
}