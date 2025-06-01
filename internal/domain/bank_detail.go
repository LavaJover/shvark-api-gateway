package domain

import "time"

type BankDetail struct {
	ID 						string
	TraderID 				string
	Country 				string
	Currency 				string
	MinAmount 				float32
	MaxAmount 				float32
	BankName 				string
	PaymentSystem 			string
	Delay					time.Duration
	Enabled 				bool
	CardNumber 				string 
	Phone 					string
	Owner 					string
	MaxOrdersSimultaneosly  int32
	MaxAmountDay			int32
	MaxAmountMonth			int32
}