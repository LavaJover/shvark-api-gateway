package request

type EditTrafficRequest struct {
	ID 						string 					`json:"id"`
	MerchantID 				*string 					`json:"merchant_id"`
	TraderID 				*string 					`json:"trader_id"`
	TraderReward 			*float64 				`json:"trader_reward"`
	TraderPriority 			*float64 				`json:"trader_priority"`
	PlatformFee 			*float64					`json:"platform_fee"`
	Enabled 				*bool 					`json:"enabled"`
	ActivityParams 			*TrafficActivityParams 	`json:"activity_params"`
	AntifraudParams 		*TrafficAntifraudParams 	`json:"antifraud_params"`
	BusinessParams 			*TrafficBusinessParams 	`json:"business_params"`
	Name					*string					`json:"name"`
}