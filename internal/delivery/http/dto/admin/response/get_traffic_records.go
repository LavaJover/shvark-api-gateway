package response

type GetTrafficResponse struct {
	TrafficRecords []Traffic `json:"traffic_records"`
}

type Traffic struct {
	ID 				string 		`json:"id"`
	MerchantID 		string 		`json:"merchant_id"`
	TraderID 		string 		`json:"trader_id"`
	TraderReward 	float64 	`json:"trader_reward"`
	TraderPriority 	float64 	`json:"trader_priority"`
	Enabled 		bool 		`json:"enabled"`
}