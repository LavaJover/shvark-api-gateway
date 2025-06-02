package request

type AddPolicyRequest struct {
	Role string `json:"role" binding:"required"`
	Object string `json:"object" binding:"required"`
	Action string `json:"action" binding:"required"`
}