package response

type AddPolicyResponse struct{
	Success bool `json:"success" binding:"required"`
}