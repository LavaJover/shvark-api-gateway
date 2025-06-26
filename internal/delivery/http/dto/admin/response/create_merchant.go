package response

type CreateMerchantResponse struct {
	MerchantID 	  string 	`json:"merchant_id"`
	AccessToken   string 	`json:"access_token"`
	WalletAddress string    `json:"wallet_address"`
}