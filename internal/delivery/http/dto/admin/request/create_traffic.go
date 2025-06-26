package request

type CreateTrafficRequest struct {
	MerchantID 		string 		`json:"merchant_id"`
	TraderID 		string 		`json:"trader_id"`
	TraderReward 	float64 	`json:"trader_reward"`
	TraderPriority 	float64 	`json:"trader_priority"`
	PlatformFee 	float64 	`json:"platform_fee"`
	Enabled 		bool    	`json:"enabled"`
}

type Traffic struct {
	ID 				string 		`json:"id"`
	MerchantID 		string 		`json:"merchant_id"`
	TraderID 		string 		`json:"trader_id"`
	TraderReward 	float64 	`json:"trader_reward"`
	TraderPriority 	float64 	`json:"trader_priority"`
	PlatformFee 	float64		`json:"platform_fee"`
	Enabled 		bool 		`json:"enabled"`
}