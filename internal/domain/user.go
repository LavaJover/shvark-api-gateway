package domain

type User struct {
	ID           string
	Login        string
	Username     string
	Password     string
	Role         string
	TwoFaSecret  string
	TwoFaEnabled bool
}

// UpdateUserData используется для обновления пользователя
type UpdateUserData struct {
	Login        string
	Username     string
	Password     string
	Role         string
	TwoFaSecret  string
	TwoFaEnabled bool
}

func NewUser(login, username, password string) *User {
	return &User{
		Login: login,
		Username: username,
		Password: password,
	}
}