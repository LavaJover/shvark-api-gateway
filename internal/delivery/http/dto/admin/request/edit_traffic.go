package request

type EditTrafficRequest struct {
	ID 						string 					`json:"id"`
	StoreID 				*string 				`json:"store_id"`
	TraderID 				*string 				`json:"trader_id"`
	TraderReward 			*float64 				`json:"trader_reward"`
	TraderPriority 			*float64 				`json:"trader_priority"`
	ActivityParams 			*TrafficActivityParams 	`json:"activity_params"`
	AntifraudParams 		*TrafficAntifraudParams `json:"antifraud_params"`
}