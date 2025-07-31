package response

// CommissionProfitResponse represents commission profit data
type CommissionProfitResponse struct {
    TraderID        string  `json:"traderId"`
    From            string  `json:"from"`
    To              string  `json:"to"`
    TotalCommission float64 `json:"totalCommission"`
    Currency        string  `json:"currency"`
}