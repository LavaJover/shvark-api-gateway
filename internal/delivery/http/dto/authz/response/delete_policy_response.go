package response

type DeletePolicyResponse struct {
	Success bool `json:"success" binding:"required"`
}