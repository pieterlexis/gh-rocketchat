package processor

import (
	"gopkg.in/go-playground/webhooks.v5/github"
)

const issueCommentCreated = `[{{ .Comment.User.Login }}]({{ .Comment.User.HTMLURL }}) [commented]({{ .Comment.HTMLURL }}) on issue.
> {{ .Comment.Body }}`

func (p *processor) handleIssueComment(ic github.IssueCommentPayload) {
	var text string
	var err error

	switch ic.Action {
	case "created":
		text, err = p.makeAndExecuteTemplate("issue_comment_created", issueCommentCreated, ic)
	case "deleted":
		text, err = p.makeAndExecuteTemplate("issue_opened", issueGenericTemplate, ic)
	default: // "edited"
		return
	}
	if err != nil {
		return
	}

	attachments, err := p.makeIssueAttachments(*ic.Issue)
	if err != nil {
		return
	}
	p.createRocketChatWebhookAndSend(
		ic.Repository.FullName,
		text,
		ic.Sender.AvatarURL,
		attachments)
}