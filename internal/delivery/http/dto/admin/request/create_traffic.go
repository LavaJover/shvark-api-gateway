package request

type CreateTrafficRequest struct {
	StoreID 				string 		`json:"store_id"`
	TraderID 				string 		`json:"trader_id"`
	TraderReward 			float64 	`json:"trader_reward"`
	TraderPriority 			float64 	`json:"trader_priority"`
	TrafficActivityParams 	TrafficActivityParams 	`json:"traffic_activity_params"`
	TrafficAntifraudParams 	TrafficAntifraudParams 	`json:"traffic_antifraud_params"`
}

type Traffic struct {
	ID 						string 					`json:"id"`
	StoreID 				string 					`json:"store_id"`
	TraderID 				string 					`json:"trader_id"`
	TraderReward 			float64 				`json:"trader_reward"`
	TraderPriority 			float64 				`json:"trader_priority"`
	TrafficActivityParams 	TrafficActivityParams 	`json:"traffic_activity_params"`
	TrafficAntifraudParams 	TrafficAntifraudParams 	`json:"traffic_antifraud_params"`
}

type TrafficActivityParams struct {
	MerchantUnlocked 	bool `json:"merchant_unlocked"`
	TraderUnlocked   	bool `json:"trader_unlocked"`
	AntifraudUnlocked 	bool `json:"antifraud_unlocked"`
	ManuallyUnlocked  	bool `json:"manually_unlocked"`
}

type TrafficAntifraudParams struct {
	AntifraudRequired bool `json:"antifraud_required"`
}