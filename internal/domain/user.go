package domain

import "time"

type User struct {
	ID 		 	string
	Login 	 	string
	Username 	string
	Password 	string
	TwoFaSecret string
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}

func NewUser(login, username, password string) *User {
	return &User{
		Login: login,
		Username: username,
		Password: password,
	}
}