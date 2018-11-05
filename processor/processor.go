package processor

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/pieterlexis/gh-rocketchat/models"
	"github.com/pieterlexis/gh-rocketchat/sender"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/webhooks.v5/github"
)

type processor struct {
	ghPayloadChan chan interface{}
	rcPayloadChan chan models.RocketChatWebhookPayload
	logPrefix     string
}

/*
  Starts a processor and a sender
*/
func RunProcessor(ghPayloadChan chan interface{}, destination string, name string) {
	rcPayloadChan := make(chan models.RocketChatWebhookPayload)
	sender.RunSender(rcPayloadChan, destination, name)
	p := processor{
		ghPayloadChan: ghPayloadChan,
		rcPayloadChan: rcPayloadChan,
		logPrefix:     fmt.Sprintf("processor(%s):", name),
	}
	go p.run()
}

func (p *processor) run() {
	for {
		select {
		case payload := <-p.ghPayloadChan:
			// TODO switch (back) to using goroutines. These were removed because the parameter was too large
			switch payload.(type) {
			case github.PingPayload:
				p.handlePing(payload.(github.PingPayload))
			case github.PullRequestPayload:
				p.handlePullRequest(payload.(github.PullRequestPayload))
			case github.PushPayload:
				p.handlePush(payload.(github.PushPayload))
			default:
				log.Warnf("%s Had an unexpected payload type: %T", p.logPrefix, payload)
			}
		}
	}
}

func (p *processor) createRocketChatWebhookAndSend(username string, text string, iconUrl string, attachments []models.RocketChatWebhookAttachment) {
	rcPayload := models.RocketChatWebhookPayload{
		Text:      text,
		LinkNames: 0, // Don't ever make rocket.chat attempt to link user and channel names
	}

	if len(iconUrl) > 0 {
		rcPayload.IconURL = iconUrl
	}

	if len(username) > 0 {
		rcPayload.UserName = username
	}

	if len(attachments) > 0 {
		rcPayload.Attachments = attachments
	}

	log.Tracef("%s Putting this rcPayload in rcPayloadChan: %+v", p.logPrefix, rcPayload)
	p.rcPayloadChan <- rcPayload
}

func (p *processor) makeAndExecuteTemplate(name string, content string, obj interface{}) (string, error) {
	t, err := p.newTemplate(name, content)
	if err != nil {
		return "", err
	}

	return p.executeTemplate(t, obj)
}

func (p *processor) newTemplate(name string, content string) (*template.Template, error) {
	t, err := template.New(name).Funcs(template.FuncMap{
		"StringsJoin": strings.Join,
		"StringSlice": func(s string, i, j int) string {
			return s[i:j]
		},
		"FirstLine": func(s string) string {
			return strings.Split(s, "\n")[0]
		},
	}).Parse(content)
	if err != nil {
		log.Warnf("%s Unable to parse template: %v", p.logPrefix, err)
		log.Debugf("%s Full template: %s", p.logPrefix, content)
		return nil, err
	}
	return t, err
}

func (p *processor) executeTemplate(template *template.Template, obj interface{}) (string, error) {
	buf := bytes.Buffer{}
	err := template.Execute(&buf, obj)

	if err != nil {
		log.Warnf("%s Could not execute template: %v", p.logPrefix, err)
		log.Tracef("%s Full payload: %+v", p.logPrefix, obj)
		return "", err
	}

	return buf.String(), nil
}
