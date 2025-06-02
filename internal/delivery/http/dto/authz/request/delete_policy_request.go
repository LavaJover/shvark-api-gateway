package request

type DeletePolicyRequest struct {
	Role string `json:"role" binding:"required"`
	Object string `json:"object" binding:"required"`
	Action string `json:"action" binding:"required"`
}