package response

type GetPayInArbitrageInfoResponse struct {
	Dispute Dispute `json:"dispute"`
}

type Dispute struct {
	DisputeID 	  string `json:"dispute_id"`
	OrderID 	  string `json:"order_id"`
	ProofUrl 	  string `json:"proof_url"`
	DisputeReason string `json:"dispute_reason"`
	DisputeStatus string `json:"dispute_status"`
}