package domain

type User struct {
	ID 		 string
	Login 	 string
	Username string
	Password string
}

func NewUser(login, username, password string) *User {
	return &User{
		Login: login,
		Username: username,
		Password: password,
	}
}