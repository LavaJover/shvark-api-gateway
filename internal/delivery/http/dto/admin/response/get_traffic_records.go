package response

type GetTrafficResponse struct {
	TrafficRecords []Traffic `json:"traffic_records"`
}

type Traffic struct {
	ID 					string
	MerchantID 			string
	TraderID 			string
	TraderRewardPercent float64
	PlatformFee			float64
	TraderPriority 		float64
	Enabled 			bool // для админов

	// Гибкие параметры
	ActivityParams 		TrafficActivityParams

	// Для антифрода
	AntifraudParams		TrafficAntifraudParams

	// Бизнес-параметры
	BusinessParams		TrafficBusinessParams
}

type TrafficActivityParams struct {
	MerchantUnlocked	bool
	TraderUnlocked		bool
	AntifraudUnlocked	bool
	ManuallyUnlocked	bool
}

type TrafficAntifraudParams struct {
	AntifraudRequired bool
}

type TrafficBusinessParams struct {
	MerchantDealsDuration string
}