package response

type GetTrafficResponse struct {
	TrafficRecords []Traffic `json:"traffic_records"`
}

type Traffic struct {
	ID 					string	`json:"id"`
	StoreID 			string `json:"store_id"`
	TraderID 			string	`json:"trader_id"`
	TraderRewardPercent float64	`json:"trader_reward_percent"`
	TraderPriority 		float64`json:"trader_priority"`

	// Гибкие параметры
	ActivityParams 		TrafficActivityParams `json:"activity_params"`

	// Для антифрода
	AntifraudParams		TrafficAntifraudParams `json:"antifraud_params"`
}

type TrafficActivityParams struct {
	MerchantUnlocked	bool `json:"merchant_unlocked"`
	TraderUnlocked		bool `json:"trader_unlocked"`
	AntifraudUnlocked	bool `json:"antifraud_unlocked"`
	ManuallyUnlocked	bool `json:"manually_unlocked"`
}

type TrafficAntifraudParams struct {
	AntifraudRequired bool `json:"antifraud_required"`
}