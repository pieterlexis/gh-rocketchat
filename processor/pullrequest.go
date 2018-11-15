package processor

import (
	"github.com/pieterlexis/gh-rocketchat/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/webhooks.v5/github"
	"strings"
)

/*
 * PullRequestEvent
 * Triggered when a pull request is assigned, unassigned, labeled, unlabeled, opened, edited, closed, reopened, or
 * synchronized. Also triggered when a pull request review is requested, or when a review request is removed.
 */

const prGenericTemplate = `[{{ .Sender.Login }}]({{ .Sender.HTMLURL }}) {{ .Action }} pull request`
const prClosedTemplate = `[{{ .Sender.Login }}]({{ .Sender.HTMLURL }}) {{ if .PullRequest.Merged }}merged{{ else }}closed{{ end }} pull request`
const prAssignedTemplate = `[{{ .Sender.Login }}]({{ .Sender.HTMLURL }}) {{ .Action }} [{{ .Assignee.Login }}]({{ .Assignee.HTMLURL }}) {{ if .Action eq "assigned" }}to{{ else }}from{{ end }} pull request`
const prLabelTemplate = `[{{ .Sender.Login }}]({{ .Sender.HTMLURL }})
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
{{- end }} pull request`
const prReviewRequestTemplate = `[{{ .Sender.Login }}]({{ .Sender.HTMLURL }}){{ if eq .Action "review_requested" }} requested a review from {{ .PullRequest.RequestedReviewer }}{{ else }} removed the review request from {{ .PullRequest.RequestedReviewer }}{{ end }} for pull request.`

const prNumberTitleTemplate = `#{{ .Number }} - {{ .PullRequest.Title }}`

func getCorrectIcon(pr github.PullRequestPayload) string {
	if pr.PullRequest.Merged {
		return ghMergedSVG
	}
	return ghPullRequestSVG
}

func (p *processor) getReviewers(pr github.PullRequestPayload) models.RocketChatWebhookField {
	var reviewers []string
	for _, r := range pr.PullRequest.RequestedReviewers {
		reviewer, _ := p.makeAndExecuteTemplate("pr-requested-reviewers", "[{{ .Login }}]({{ .HTMLURL }})", r)
		reviewers = append(reviewers, reviewer)
	}
	return models.RocketChatWebhookField{
		Title: "requested reviewers",
		Value: strings.Join(reviewers, ", "),
	}
}

func (p *processor) getLabels(pr github.PullRequestPayload) models.RocketChatWebhookField {
	var labels []string
	for _, l := range pr.PullRequest.Labels {
		labels = append(labels, l.Name)
	}
	return models.RocketChatWebhookField{
		Title: "labels",
		Value: strings.Join(labels, ", "),
		// Short: len(labels) <= 40,
	}
}

func (p *processor) getAssignees(pr github.PullRequestPayload) models.RocketChatWebhookField {
	var assignees []string
	for _, a := range pr.PullRequest.Assignees {
		assignee, _ := p.makeAndExecuteTemplate("pr-assignees", "[{{ .Login }}]({{ .HTMLURL }})", a)
		assignees = append(assignees, assignee)
	}
	return models.RocketChatWebhookField{
		Title: "assignees",
		Value: strings.Join(assignees, ", "),
	}
}

func (p *processor) getPullRequestBody(pr github.PullRequestPayload) models.RocketChatWebhookField {
	return models.RocketChatWebhookField{
		Title: "body",
		Value: pr.PullRequest.Body,
		// Short: len(pr.PullRequest.Body) >= 512,
	}
}

func (p *processor) getPullRequestMilestone(pr github.PullRequestPayload) models.RocketChatWebhookField {
	milestone, _ := p.makeAndExecuteTemplate("issue-milestone", "[{{ .PullRequest.Milestone.Title }}]({{ .PullRequest.Milestone.HTMLURL }})", pr)
	return models.RocketChatWebhookField{
		Title: "milestone",
		Value: milestone,
	}
}

func (p *processor) makePullRequestAttachments(pr github.PullRequestPayload) ([]models.RocketChatWebhookAttachment, error) {
	var attachments []models.RocketChatWebhookAttachment
	var fields []models.RocketChatWebhookField

	prTitleAndName, err := p.makeAndExecuteTemplate("pr_title_name", prNumberTitleTemplate, pr)
	if err != nil {
		return attachments, err
	}

	fields = append(fields, models.RocketChatWebhookField{
		Title: "opened",
		Value: pr.PullRequest.User.Login,
	})

	if len(pr.PullRequest.Labels) > 0 {
		fields = append(fields, p.getLabels(pr))
	}

	if len(pr.PullRequest.Assignees) > 0 {
		fields = append(fields, p.getAssignees(pr))
	}

	if pr.PullRequest.Milestone != nil {
		fields = append(fields, p.getPullRequestMilestone(pr))
	}

	if len(pr.PullRequest.RequestedReviewers) > 0 {
		fields = append(fields, p.getReviewers(pr))
	}

	fields = append(fields, p.getPullRequestBody(pr))

	attachments = append(attachments, models.RocketChatWebhookAttachment{
		Title: "Details",
		Collapsed: true,
		AuthorIcon: svgInlinePrefix + getCorrectIcon(pr),
		AuthorName: prTitleAndName,
		AuthorLink: pr.PullRequest.HTMLURL,
		Fields: fields,
	})

	return attachments, nil
}

func (p *processor) handlePullRequest(pr github.PullRequestPayload) {
	var text string
	var err error

	switch pr.Action {
	case "assigned", "unassigned":
		text, err = p.makeAndExecuteTemplate("pr_assigned", prAssignedTemplate, pr)
	case "closed":
		text, err = p.makeAndExecuteTemplate("pr_closed", prClosedTemplate, pr)
	case "labeled", "unlabeled":
		text, err = p.makeAndExecuteTemplate("pr_labeled", prLabelTemplate, pr)
	case "review_requested", "review_requested_removed":
		text, err = p.makeAndExecuteTemplate("pr_review_request", prReviewRequestTemplate, pr)
	case "opened", "reopened":
		text, err = p.makeAndExecuteTemplate("pr_opened", prGenericTemplate, pr)
	default:
		logrus.Infof("%s Unhandled Pull Request action '%s'", p.logPrefix, pr.Action)
		return
	}

	if err != nil {
		return
	}

	attachments, err := p.makePullRequestAttachments(pr)
	if err != nil {
		return
	}

	p.createRocketChatWebhookAndSend(
		pr.Repository.FullName,
		text,
		pr.PullRequest.User.AvatarURL,
		attachments)
}
