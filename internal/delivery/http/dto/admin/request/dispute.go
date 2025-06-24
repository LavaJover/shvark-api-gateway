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
	OrderID 		string `json:"order_id"`
	ProofUrl 		string `json:"proof_url"`
	DisputeReason 	string `json:"dispute_reason"`
	Ttl 			string `json:"ttl"`
}

type GetDisputeInfoRequest struct {
	DisputeID string `json:"dispute_id"`
}