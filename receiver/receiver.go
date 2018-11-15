package receiver

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/webhooks.v5/github"
	"net/http"
)

type receiver struct {
	hook      *github.Webhook
	logPrefix string
}

func NewReceiver(name string, secret string) (*receiver, error) {
	hook, err := github.New(github.Options.Secret(secret))
	if err != nil {
		return nil, err
	}

	return &receiver{
		hook:      hook,
		logPrefix: fmt.Sprintf("receiver(%s):", name),
	}, nil
}

func (r *receiver) Handle(ghPayloadChan chan interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ghPayload, err := r.hook.Parse(req,
			github.IssuesEvent,
			github.PingEvent,
			github.PullRequestEvent,
			github.IssueCommentEvent,
			github.PushEvent)

		if err != nil {
			if err == github.ErrEventNotFound {
				return
			}
			log.Warnf("%s Could not parse hook: %s", r.logPrefix, err)
			log.Tracef("%s Request content: %+v", r.logPrefix, req)
			return
		}

		log.Debugf("%s Had event payload: %T", r.logPrefix, ghPayload)
		log.Tracef("%s Event content: %+v", r.logPrefix, ghPayload)

		ghPayloadChan <- ghPayload
	})
}
