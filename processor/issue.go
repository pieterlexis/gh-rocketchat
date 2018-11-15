package processor

import (
	"github.com/pieterlexis/gh-rocketchat/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/webhooks.v5/github"
	"strings"
)

const issueGenericTemplate = `[{{ .Sender.Login }}]({{ .Sender.HTMLURL }}) {{ .Action }} issue`
const issueAssignedTemplate = `[{{ .Sender.Login }}]({{ .Sender.HTMLURL }}) {{ .Action }} [{{ .Assignee.Login }}]({{ .Assignee.HTMLURL }}) {{ if .Action eq "assigned" }}to{{ else }}from{{ end }} issue`
const issueLabelTemplate = `[{{ .Sender.Login }}]({{ .Sender.HTMLURL }})
{{- if eq .Action "labeled" }}
	{{- " added " }}
{{- else }}
	{{ " removed " }}
{{- end }}
{{- "label " }}
{{- .Label.Name }}
{{- if eq .Action "labeled" }}
{{- "to " }}
{{- else }}
{{ "from " }}
{{- end }} issue`

const issueNumberTitleTemplate = `#{{ .Number }} - {{ .Title }}`

const issueMilestonedTemplate = `[{{ .Sender.Login }}]({{ .Sender.HTMLURL }})
{{- if eq .Action "milestoned" }}
	{{- " added " }}
{{- else }}
	{{ " removed " }}
{{- end }}
{{- "the milestone " }}
{{- .Issue.Milestone.Title }}
{{- if eq .Action "milestoned" }}
{{- " to " }}
{{- else }}
{{ " from " }}
{{- end }}issue`

func (p *processor) getIssueBody(issue github.Issue) models.RocketChatWebhookField {
	return models.RocketChatWebhookField{
		Title: "body",
		Value: issue.Body,
		// Short: len(issue.PullRequest.Body) >= 512,
	}
}

func (p *processor) getIssueLabels(issue github.Issue) models.RocketChatWebhookField {
	var labels []string
	for _, l := range issue.Labels {
		labels = append(labels, l.Name)
	}
	return models.RocketChatWebhookField{
		Title: "labels",
		Value: strings.Join(labels, ", "),
		// Short: len(labels) <= 40,
	}
}

func (p *processor) getIssueAssignees(issue github.Issue) models.RocketChatWebhookField {
	var assignees []string
	for _, a := range issue.Assignees {
		assignee, _ := p.makeAndExecuteTemplate("issue-assignees", "[{{ .Login }}]({{ .HTMLURL }})", a)
		assignees = append(assignees, assignee)
	}
	return models.RocketChatWebhookField{
		Title: "assignees",
		Value: strings.Join(assignees, ", "),
	}
}

func (p *processor) getIssueMilestone(issue github.Issue) models.RocketChatWebhookField {
	milestone, _ := p.makeAndExecuteTemplate("issue-milestone", "[{{ .Milestone.Title }}]({{ .Milestone.HTMLURL }})", issue)
	return models.RocketChatWebhookField{
		Title: "milestone",
		Value: milestone,
	}
}

func (p *processor) makeIssueAttachments(issue github.Issue) ([]models.RocketChatWebhookAttachment, error) {
	// TODO roll up with makePullRequestAttachments

	var attachments []models.RocketChatWebhookAttachment
	var fields []models.RocketChatWebhookField

	issueTitleAndName, err := p.makeAndExecuteTemplate("issue_title_name", issueNumberTitleTemplate, issue)
	if err != nil {
		return attachments, err
	}

	fields = append(fields, models.RocketChatWebhookField{
		Title: "opened",
		Value: issue.User.Login,
	})

	if len(issue.Labels) > 0 {
		fields = append(fields, p.getIssueLabels(issue))
	}

	if issue.Assignees != nil {
		fields = append(fields, p.getIssueAssignees(issue))
	}

	if issue.Milestone != nil {
		fields = append(fields, p.getIssueMilestone(issue))
	}

	fields = append(fields, p.getIssueBody(issue))

	attachments = append(attachments, models.RocketChatWebhookAttachment{
		Title:      "Details",
		Collapsed:  true,
		AuthorIcon: svgInlinePrefix + ghIssueSVG,
		AuthorName: issueTitleAndName,
		AuthorLink: issue.HTMLURL,
		Fields:     fields,
	})

	return attachments, nil
}

func (p *processor) handleIssue(issue github.IssuesPayload) {
	var text string
	var err error

	switch issue.Action {
	case "assigned", "unassigned":
		text, err = p.makeAndExecuteTemplate("issue_assigned", issueAssignedTemplate, issue)
	case "labeled", "unlabeled":
		text, err = p.makeAndExecuteTemplate("issue_labeled", issueLabelTemplate, issue)
	case "opened", "reopened", "closed":
		text, err = p.makeAndExecuteTemplate("issue_opened", issueGenericTemplate, issue)
	case "milestoned", "unmilestoned":
		text, err = p.makeAndExecuteTemplate("issue_milestoned", issueMilestonedTemplate, issue)
	default:
		logrus.Infof("%s Unhandled Issue action '%s'", p.logPrefix, issue.Action)
		return
	}

	if err != nil {
		return
	}

	attachments, err := p.makeIssueAttachments(*issue.Issue)
	if err != nil {
		return
	}

	p.createRocketChatWebhookAndSend(
		issue.Repository.FullName,
		text,
		issue.Issue.User.AvatarURL,
		attachments)
}
