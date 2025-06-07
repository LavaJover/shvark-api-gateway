package response

type ReleaseResponse struct {
	Released float64 `json:"released"`
	Reward float64 `json:"reward"`
}

type ReleaseErrorResponse struct {
	Error string `json:"error"`
}