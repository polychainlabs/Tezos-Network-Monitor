package alert

import (
	"log"
	"os"
	"time"

	"github.com/nlopes/slack"
)

// PostSlack message
func PostSlack(msg *slack.WebhookMessage) {
	// Set Defaults
	if len(msg.Channel) == 0 {
		msg.Channel = os.Getenv("SLACK_CHANNEL")
	}
	// Throttle
	if hasAlreadyAlerted(msg.Text+msg.Channel, 10*time.Minute) {
		log.Println("Would have sent alert text: ", msg.Text)
		log.Println("Has already alerted.  Not posting again.")
		return
	}
	// Post
	err := slack.PostWebhook(os.Getenv("SLACK_URL"), msg)
	if err != nil {
		log.Println("Error posting slack webhook: ", err)
	}
}
