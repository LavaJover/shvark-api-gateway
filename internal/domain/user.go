package domain

type User struct {
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