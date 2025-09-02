package request

type GetBankDetailsQuery struct {
	TraderID 		*string `form:"trader_id"`
	BankCode 		*string `form:"bank_code"`
	Enabled  		*bool   `form:"enabled"`
	PaymentSystem 	*string `form:"payment_system"`
	BankDetailID 	*string  `form:"bank_detail_id"`
	Page 			int 	`form:"page"`
	Limit			int 	`form:"limit"`
}