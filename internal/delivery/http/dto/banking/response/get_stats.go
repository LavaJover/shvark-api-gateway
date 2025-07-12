package response

type BankDetailStat struct {
	BankDetailID   	   string  	`json:"bank_detail_id"`
	CurrentCountToday  int     	`json:"current_count_today"`
	CurrentCountMonth  int     	`json:"current_count_month"`
	CurrentAmountToday float64	`json:"current_amount_today"`
	CurrentAmountMonth float64	`json:"current_amount_month"`
}

type GetBankDetailsStatsResponse struct {
	Stats []BankDetailStat `json:"stats"`
}