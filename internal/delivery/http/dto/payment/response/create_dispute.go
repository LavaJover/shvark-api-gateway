package response

type CreateDisputeResponse struct {
	DisputeID string `json:"dispute_id"`
	ExpiresAt string `json:"expires_at"`
}