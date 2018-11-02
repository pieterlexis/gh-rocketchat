package processor

import (
	"strings"

	"github.com/pieterlexis/gh-rocketchat/models"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const pushTemplate string = `[{{ .Sender.Login }}]({{ .Sender.HTMLURL }}) {{ if .Forced }}*force* {{ end }}pushed {{ len .Commits }} commits to [{{ .Ref }}]({{ .Compare }}) in [{{ .Repository.FullName }}]({{ .Repository.URL }})`

const pushCommitTemplate string = `[{{ .ID }}]({{ .URL }}): `

func (p *processor) handlePush(push github.PushPayload) {
	text, err := p.makeAndExecuteTemplate("push", pushTemplate, push)
	if err != nil {
		return
	}

	var commits = make([]models.RocketChatWebhookField, len(push.Commits))

	for _, commit := range push.Commits {
		msg := strings.Split(commit.Message, "\n")[0]
		prefix, _ := p.makeAndExecuteTemplate("pushCommit", pushCommitTemplate, commit)
		commits = append(commits, models.RocketChatWebhookField{
			Title: commit.Author.Name,
			Value: prefix + msg,
		})
	}
	attachments := []models.RocketChatWebhookAttachment{{
		Title:  "Commits: ",
		Fields: commits,
	}}

	p.createRocketChatWebhookAndSend(
		push.Sender.Login,
		text,
		push.Sender.AvatarURL,
		attachments)
}
