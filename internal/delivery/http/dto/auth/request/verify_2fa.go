package request

type Verify2FARequest struct{
	Code   string `json:"code"` 
}