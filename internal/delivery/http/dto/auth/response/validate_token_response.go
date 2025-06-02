package response

type ValidateTokenResponse struct {
	Valid  bool 	`json:"valid"`
	UserId string 	`json:"user_id" example:"1d6ab366-4fca-4bcc-972d-875c35ea939a"`
}