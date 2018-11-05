package models

/* See https://rocket.chat/docs/developer-guides/rest-api/chat/postmessage/

 */

type RocketChatWebhookField struct {
	Short bool   `json:"short"`
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
}

type RocketChatWebhookAttachment struct {
	AudioUrl          string                   `json:"audio_url,omitempty"`
	AuthorIcon        string                   `json:"author_icon,omitempty"`
	AuthorLink        string                   `json:"author_link,omitempty"`
	AuthorName        string                   `json:"author_name,omitempty"`
	Collapsed         bool                     `json:"collapsed"`
	Color             string                   `json:"color,omitempty"`
	Fields            []RocketChatWebhookField `json:"fields,omitempty"`
	ImageUrl          string                   `json:"image_url,omitempty"`
	MessageLink       string                   `json:"message_link,omitempty"`
	Text              string                   `json:"text,omitempty"`
	ThumbUrl          string                   `json:"thumb_url,omitempty"`
	Title             string                   `json:"title,omitempty"`
	TitleLink         string                   `json:"title_link,omitempty"`
	TitleLinkDownload bool                     `json:"title_link_download"`
	Ts                string                   `json:"ts,omitempty"`
	VideoUrl          string                   `json:"video_url,omitempty"`
}

type RocketChatWebhookPayload struct {
	Attachments []RocketChatWebhookAttachment `json:"attachments,omitempty"`
	Channel     string                        `json:"channel,omitempty"`
	IconURL     string                        `json:"icon_url,omitempty"`
	LinkNames   int                           `json:"link_names"`
	Text        string                        `json:"text,omitempty"`
	UserName    string                        `json:"username,omitempty"`
}
