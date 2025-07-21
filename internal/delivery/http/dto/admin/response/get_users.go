package response

type GetUsersResponse struct {
	Users []User `json:"users"`
}

type User struct {
	ID 		 string `json:"id"`
	Username string	`json:"username"`
	Login 	 string `json:"login"`
	Role 	 string `json:"role"`
}