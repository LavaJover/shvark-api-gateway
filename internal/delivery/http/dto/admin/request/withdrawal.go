package request

type SetWithdrawalRulesRequest struct {
	UserID 			string 	`json:"user_id"`
	FixedFee 		float64 `json:"fixed_fee"`
	MinAmount 		float64 `json:"min_amount"`
	CooldownSeconds int64 	`json:"cooldown_seconds"`
}