package response

type UserResponse struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	TwoFaEnabled bool   `json:"two_fa_enabled"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type UserDetailResponse struct {
	UserResponse
	TwoFaSecret string `json:"two_fa_secret,omitempty"`
}

type UserWithPasswordResponse struct {
	UserResponse
	Password string `json:"password,omitempty"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

type UpdateUserResponse struct {
	User UserResponse `json:"user"`
}

type DeleteUserResponse struct {
	Success bool `json:"success"`
}

type ChangePasswordResponse struct {
	Success bool `json:"success"`
}

type GetUsersResponse struct {
	TotalPages int64           `json:"total_pages"`
	TotalCount int64           `json:"total_count"`
	Users      []UserResponse  `json:"users"`
}

type TwoFaSecretResponse struct {
	TwoFaSecret string `json:"two_fa_secret"`
}

type SetTwoFaEnabledResponse struct {
	Success bool `json:"success"`
}