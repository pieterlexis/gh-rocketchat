package processor

import (
	"github.com/pieterlexis/gh-rocketchat/models"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const pingTemplate string = `:thumbsup: Received PING with ID '{{ .HookID }}' successfully!`

func (p *processor) handlePing(ping github.PingPayload) {
	text, err := p.makeAndExecuteTemplate("ping", pingTemplate, ping)
	if err != nil {
		return
	}
	p.createRocketChatWebhookAndSend(
		"GitHub",
		text,
		"https://assets-cdn.github.com/images/modules/logos_page/Octocat.png",
		[]models.RocketChatWebhookAttachment{})
}
