package request

type UpdateUserRequest struct {
	Login        *string `json:"login,omitempty" binding:"omitempty,min=3,max=50"`
	Username     *string `json:"username,omitempty" binding:"omitempty,min=2,max=100"`
	Password     *string `json:"password,omitempty" binding:"omitempty,min=6,max=100"`
	TwoFaSecret  *string `json:"two_fa_secret,omitempty"`
	TwoFaEnabled *bool   `json:"two_fa_enabled,omitempty"`
	Role         *string `json:"role,omitempty" binding:"omitempty,oneof=TRADER MERCHANT ADMIN TEAM_LEAD"`
}