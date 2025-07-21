package response

import "time"

type SetWithdrawalRulesResponse struct {
	Success bool `json:"success"`
	Rule 	Rule `json:"rule"`
}

type GetWithdrawalRulesResponse struct {
	Rule Rule `json:"rule"`
}

type DeleteWithdrawalRulesResponse struct {

}

type Rule struct {
	ID 				int64 	  `json:"id"`
	TraderID 		string 	  `json:"trader_id"`
	FixedFee 		float64   `json:"fixed_fee"`
	MinAmount 		float64   `json:"min_amount"`
	CooldownSeconds int64 	  `json:"cooldown_seconds"`
	UpdatedAt 		time.Time `json:"updated_at"`
	CreatedAt 		time.Time `json:"created_at"`
}