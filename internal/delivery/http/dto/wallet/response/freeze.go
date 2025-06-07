package response

type FreezeResponse struct {
	Frozen float64 `json:"frozen"`
}

type FreezeErrorResponse struct {
	Error string `json:"error" example:"Insufficient balance"`
}