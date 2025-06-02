package response

type GetProfileByIDResponse struct {
	ProfileID string `json:"profile_id"`
	AvatarURL string `json:"avatar_url"`
	TgLink 	  string `json:"tg_link"`
	UserID    string `json:"user_id"`
}