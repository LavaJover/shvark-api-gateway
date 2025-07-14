package response

type GetUsersResponse struct {
	Users []User `json:"users"`
}

type User struct {
	Username string	`json:"username"`
	Login 	 string `json:"login"`
	Role 	 string `json:"role"`
}