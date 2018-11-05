package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pieterlexis/gh-rocketchat/models"
	log "github.com/sirupsen/logrus"
)

type sender struct {
	rcPayloadChan chan models.RocketChatWebhookPayload
	delayedQueue  chan delayedMsg
	destination   string
	logPrefix     string
}

func RunSender(rcPayloadChan chan models.RocketChatWebhookPayload, destination string, name string) {
	s := sender{
		rcPayloadChan: rcPayloadChan,
		delayedQueue:  make(chan delayedMsg),
		destination:   destination,
		logPrefix:     fmt.Sprintf("sender(%s):", name),
	}
	go s.run()
}

func (s *sender) run() {
	delayedTicker := time.NewTicker(5 * time.Second)
	defer delayedTicker.Stop()
	log.Infof("%s started!", s.logPrefix)

	for {
		select {
		case msg := <-s.rcPayloadChan:
			log.Tracef("%s received payload on rcPayloadChan: %+v", s.logPrefix, msg)
			go s.sendToWebhook(msg)
		case <-delayedTicker.C:
			log.Tracef("%s delayTicker ticked", s.logPrefix)
			go s.handleDelayedMsgs()
		}
	}
}

func (s *sender) sendToWebhook(payload models.RocketChatWebhookPayload) {
	jsonValue, err := json.Marshal(payload)
	if err != nil {
		log.Warnf("%s Could not marshall message to JSON: %+v", s.logPrefix, err)
		log.Debugf("%s Full payload value: %+v", s.logPrefix, payload)
	}

	log.Debugf("%s Sending webhook to %s", s.logPrefix, s.destination)
	log.Tracef("%s with this content: %s", s.logPrefix, jsonValue)

	_, err = http.Post(s.destination, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Warnf("%s Had error while sending webhook to Rocket.Chat: %v", s.logPrefix, err)
		// TODO Handle unreachable errors by pushing the message with info to the delayedMsg chan
		return
	}
}

func (s *sender) handleDelayedMsgs() {
	// TODO try to send everything in the delayQueue if they have been in there for longer than X seconds (lastAttempt)
	// TODO If the hooks cannot be delivered, update lastAttempt and stick it back in the queue.
	// TODO Unless the originalTime is older than Y seconds
}
