package models

type RocketChatWebhookField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type RocketChatWebhookAttachment struct {
	Color      string                   `json:"color"`
	AuthorName string                   `json:"author_name"`
	AuthorLink string                   `json:"author_link"`
	AuthorIcon string                   `json:"author_icon"`
	Title      string                   `json:"title"`
	TitleLink  string                   `json:"title_link"`
	Text       string                   `json:"text"`
	Fields     []RocketChatWebhookField `json:"fields"`
	ImageUrl   string                   `json:"image_url"`
	ThumbUrl   string                   `json:"thumb_url"`
}

type RocketChatWebhookPayload struct {
	UserName    string                        `json:"username"`
	LinkNames   int                           `json:"link_names"`
	Channel     string                        `json:"channel"`
	IconURL     string                        `json:"icon_url"`
	Text        string                        `json:"text"`
	Attachments []RocketChatWebhookAttachment `json:"attachments"`
}
