package response

type GetUserByIDResponse struct {
	UserID string `json:"user_id"`
	Login string `json:"login"`
	Username string `json:"username"`
	Password string `json:"password"`
}