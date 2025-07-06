package request

type CreateDisputeRequest struct {
	ProofUrl string 	`json:"proof_url"`
	Reason   string  	`json:"reason"`
}