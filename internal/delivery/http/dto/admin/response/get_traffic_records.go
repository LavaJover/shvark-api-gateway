package response

type GetTrafficResponse struct {
	TrafficRecords []Traffic `json:"traffic_records"`
}

type Traffic struct {
	ID 					string	`json:"id"`
	MerchantID 			string `json:"merchant_id"`
	TraderID 			string	`json:"trader_id"`
	TraderRewardPercent float64	`json:"trader_reward_percent"`
	PlatformFee			float64	`json:"platform_fee"`
	TraderPriority 		float64`json:"trader_priority"`
	Enabled 			bool // для админов `json:"enabled"`
	Name				string `json:"name"`

	// Гибкие параметры
	ActivityParams 		TrafficActivityParams `json:"activity_params"`

	// Для антифрода
	AntifraudParams		TrafficAntifraudParams `json:"antifraud_params"`

	// Бизнес-параметры
	BusinessParams		TrafficBusinessParams `json:"business_params"`
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

type TrafficBusinessParams struct {
	MerchantDealsDuration string `json:"merchant_deals_duration"`
}