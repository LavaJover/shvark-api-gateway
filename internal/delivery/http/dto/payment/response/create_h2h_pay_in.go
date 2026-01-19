package response

type CreateH2HPayInResponse struct {
	OrderID 		string  	   `json:"order_id"`
	AmountFiat 		float64		   `json:"amount_fiat"`
	AmountCrypto 	float64		   `json:"amount_crypto"`
	Currency 		string 		   `json:"currency"`
	PaymentSystem	string		   `json:"payment_system"`
	Status 			string 		   `json:"status"`
	MerchantOrderID	string		   `json:"merchant_order_id"`
	CallbackURL		string		   `json:"callback_url"`
	PaymentDetails 	PaymentDetails `json:"payment_details"`
	ExpiresAt 		int64 		   `json:"expires_at"`
	TPayLink 		string 		   `json:"tpay_link"`
	Recalculated	bool		   `json:"recalculated"`
	CryptoRubRate 	float64	   `json:"crypto_rub_rate"`
	DeeplinkHTML    string			`json:"deeplink_html"`
	DeeplinkRedirect string			`json:"deeplink_redirect"`
	MerchantIncome	float64		`json:"merchant_income"`
	UsdRate			float64		`json:"usd_rate"`
}

type NoBankDetailsErrorResponse struct {
	Error string `json:"error"`
}

type BadRequestErrorResponse struct {
	Error string `json:"error"`
}