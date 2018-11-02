package sender

import (
	"time"

	"github.com/pieterlexis/gh-rocketchat/models"
)

type delayedMsg struct {
	payload      models.RocketChatWebhookPayload
	originalTime time.Time
	lastAttempt  time.Time
}
