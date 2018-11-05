package processor

import (
	"strings"

	"github.com/pieterlexis/gh-rocketchat/models"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const pushTemplate string = `[{{ .Sender.Login }}]({{ .Sender.HTMLURL }}) {{ if .Forced }}*force* {{ end }}pushed {{ len .Commits }} commits to [{{ .Ref }}]({{ .Compare }}) in [{{ .Repository.FullName }}]({{ .Repository.URL }})`
const pushCommitTemplate string = `[{{ StringSlice .ID 0 8 }}]({{ .URL }}): {{ FirstLine .Message }} ({{ .Author.Name }})`

func (p *processor) handlePush(push github.PushPayload) {
	text, err := p.makeAndExecuteTemplate("push", pushTemplate, push)
	if err != nil {
		return
	}

	var commits []string
	for _, commit := range push.Commits {
		msg, _ := p.makeAndExecuteTemplate("pushCommitPrefix", pushCommitTemplate, commit)
		commits = append(commits, msg)
	}

	var commitsField = models.RocketChatWebhookField{
		Title: "commits",
		Value: strings.Join(commits, "\n"),
	}

	attachments := []models.RocketChatWebhookAttachment{{
		Title:     "Details",
		Collapsed: true,
		Fields: []models.RocketChatWebhookField{
			commitsField,
		},
	}}

	p.createRocketChatWebhookAndSend(
		push.Repository.FullName,
		text,
		push.Sender.AvatarURL,
		attachments)
}
