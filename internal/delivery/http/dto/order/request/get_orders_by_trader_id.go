package request

type GetOrdersByTraderIDRequest struct {
	TraderID string `json:"trader_id"`
}

type OrderQueryParams struct {
	Page 		int64 		`form:"page" binding:"min=1"`
	Limit 		int64 		`form:"limit" binding:"min=1,max=100"`
	SortBy 		string 		`form:"sort_by"`
	SortOrder	string 		`form:"sort_order"`
	Status 		[]string 	`form:"status"`
	MinAmount 	float64		`form:"min_amount"`
	MaxAmount	float64     `form:"max_amount"`
	DateFrom 	string 		`form:"date_from"`
	DateTo		string   	`form:"date_to"`
	Currency	string 		`form:"currency"`
}