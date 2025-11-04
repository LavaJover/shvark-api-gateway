package request

import "time"

type GetOrdersParams struct {
	DealID				*string 	`form:"dealId"`
	Type 				*string 	`form:"type"`
	Status 				*string 	`form:"status"`
	TimeOpeningStart 	*time.Time 	`form:"timeOpeningStart" time_format:"2006-01-02T15:04:05Z"`
	TimeOpeningEnd 		*time.Time 	`form:"timeOpeningEnd" time_format:"2006-01-02T15:04:05Z"`
	AmountMin 			*float64	`form:"amountMin"`
	AmountMax			*float64	`form:"amountMax"`
	Page 				*int		`form:"page"`
	Size 				*int		`form:"size"`
	Sort 				*string		`form:"sort"`
}