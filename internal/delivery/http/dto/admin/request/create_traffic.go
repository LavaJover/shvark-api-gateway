package request

import "time"

type CreateTrafficRequest struct {
	MerchantID 				string 		`json:"merchant_id"`
	TraderID 				string 		`json:"trader_id"`
	TraderReward 			float64 	`json:"trader_reward"`
	TraderPriority 			float64 	`json:"trader_priority"`
	PlatformFee 			float64 	`json:"platform_fee"`
	Enabled 				bool    	`json:"enabled"`
	TrafficActivityParams 	TrafficActivityParams 	`json:"traffic_activity_params"`
	TrafficAntifraudParams 	TrafficAntifraudParams 	`json:"traffic_antifraud_params"`
	TrafficBusinessParams 	TrafficBusinessParams 	`json:"traffic_business_params"`
}

type Traffic struct {
	ID 						string 					`json:"id"`
	MerchantID 				string 					`json:"merchant_id"`
	TraderID 				string 					`json:"trader_id"`
	TraderReward 			float64 				`json:"trader_reward"`
	TraderPriority 			float64 				`json:"trader_priority"`
	PlatformFee 			float64					`json:"platform_fee"`
	Enabled 				bool 					`json:"enabled"`
	TrafficActivityParams 	TrafficActivityParams 	`json:"traffic_activity_params"`
	TrafficAntifraudParams 	TrafficAntifraudParams 	`json:"traffic_antifraud_params"`
	TrafficBusinessParams 	TrafficBusinessParams 	`json:"traffic_business_params"`
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

type TrafficBusinessParams struct {
	MerchantDealsDuration time.Duration `json:"merchant_deals_duration"`
}