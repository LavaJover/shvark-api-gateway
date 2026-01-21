package response

type CreateTrafficResponse struct {
	Traffic Traffic `json:"traffic"`
	Message string `json:"message"`
}