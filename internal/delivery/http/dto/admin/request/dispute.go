package request

type AcceptDisputeRequest struct {
	DisputeID string `json:"dispute_id"`
}

type RejectDisputeRequest struct {
	DisputeID string `json:"dispute_id"`
}

type FreezeDisputeRequest struct {
	DisputeID string `json:"dispute_id"`
}

type CreateDisputeRequest struct{
	OrderID 			string `json:"order_id"`
	ProofUrl 			string `json:"proof_url"`
	DisputeReason 		string `json:"dispute_reason"`
	Ttl 				string `json:"ttl"`
	DisputeAmountFiat 	float64 `json:"dispute_amount_fiat"`
}

type GetDisputeInfoRequest struct {
	DisputeID string `json:"dispute_id"`
}

type GetOrderDisputesQuery struct {
	TraderID   *string `form:"traderId"`
	MerchantID *string `form:"merchantId"`
	OrderID    *string `form:"orderId"`
	DisputeID  *string `form:"disputeId"`
	Status 	   *string `form:"status"`
	Page 	   int		`form:"page"`
	Limit 	   int 		`form:"limit"`  	
}