package merchant

type CreatePayInRequest struct {
	IsSbp 		bool 	`json:"isSbp"`
	Amount 		float64 `json:"amount"`
	Currency 	string 	`json:"currency"`
	Issuer 		string 	`json:"issuer"`
	NspkCode 	string 	`json:"nspkCode"`
	CallbackUrl string 	`json:"callbackUrl"`
	IternalID 	string 	`json:"iternalId"`
}

type CreatePayInResponse struct {
	OrderID 		 string  `json:"orderId"`
	CardNumber 		 string  `json:"cardNumber"`
	PhoneNumber 	 string  `json:"phoneNumber"`
	HolderName  	 string  `json:"holderName"`
	Issuer 			 string  `json:"issuer"`
	NspkCode 		 string  `json:"nspkCode"`
	Amount 			 float64 `json:"amount"`
	AmountByCurrency float64 `json:"amountByCurrency"`
	CurrencyRate 	 float64 `json:"currencyRate"`
	TimeExpires 	 int64  `json:"timeExpires"`
}