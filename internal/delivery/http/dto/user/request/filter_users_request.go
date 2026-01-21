package request

type FilterUsersRequest struct {
	Page           int64    `form:"page" binding:"omitempty,min=1"`
	Limit          int64    `form:"limit" binding:"omitempty,min=1,max=100"`
	LoginFilter    string   `form:"login_filter" binding:"omitempty,max=50"`
	UsernameFilter string   `form:"username_filter" binding:"omitempty,max=100"`
	Roles          []string `form:"roles" binding:"omitempty,dive,oneof=TRADER MERCHANT ADMIN TEAM_LEAD"`
	SortBy         string   `form:"sort_by" binding:"omitempty,oneof=login username created_at role"`
	SortDesc       bool     `form:"sort_desc"`
}