package request

type CreateDisputeRequest struct {
	OrderID  string 	`json:"order_id"`
	ProofUrl string 	`json:"proof_url"`
}