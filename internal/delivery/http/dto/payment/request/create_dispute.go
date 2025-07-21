package request

type CreateDisputeRequest struct {
	ProofUrl string 	`json:"proof_url"`
	Reason   string  	`json:"reason"`
	AmountFiat	float64	`json:"amount_fiat"`
}