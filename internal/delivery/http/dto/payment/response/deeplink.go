package response

type DeeplinkResponse struct {
	HTMLContent string `json:"html_content,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
	DeeplinkURL string `json:"deeplink_url,omitempty"`
}