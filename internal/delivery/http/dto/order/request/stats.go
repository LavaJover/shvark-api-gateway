package request

import "time"

type GetOrderStatsRequest struct {
	TraderID string 	`json:"trader_id"`
	DateFrom time.Time 	`json:"date_from"`
	DateTo 	 time.Time 	`json:"date_to"`
}