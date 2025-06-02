package request

type UpdateBankDetailRequest struct {
	BankDetail BankDetail `json:"bank_detail"`
}

type BankDetail struct {
	ID 						string	`json:"id"`
	TraderID 				string	`json:"trader_id"`
	Country 				string	`json:"country"`
	Currency 				string	`json:"currency"`
	MinAmount 				float32	`json:"min_amount"`
	MaxAmount 				float32	`json:"max_amount"`
	BankName 				string	`json:"bank_name"`
	PaymentSystem 			string	`json:"payment_system"`
	Delay					string	`json:"delay"`
	Enabled 				bool	`json:"enabled"`
	CardNumber 				string  `json:"card_number"`
	Phone 					string	`json:"phone"`
	Owner 					string	`json:"owner"`
	MaxOrdersSimultaneosly  int32	`json:"max_orders_simultaneosly"`
	MaxAmountDay			int32	`json:"max_amount_day"`
	MaxAmountMonth			int32	`json:"max_amount_month"`
}