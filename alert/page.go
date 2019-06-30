package alert

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/nlopes/slack"
)

// Page your team
func Page(title string, body string) {
	if hasAlreadyAlerted(title, 20*time.Minute) {
		return
	}

	// If we're going to page, always slack alert too
	PostSlack(&slack.WebhookMessage{
		Text: fmt.Sprintf("*Paging* with title: _%v_", title),
	})

	//Page
	pd := pagerduty.NewClient(os.Getenv("PD_TOKEN"))
	_, err := pd.CreateIncident(os.Getenv("PD_USER"), &pagerduty.CreateIncident{
		Incident: pagerduty.CreateIncidentOptions{
			Type:  "incident",
			Title: title,
			Service: pagerduty.APIReference{
				ID:   os.Getenv("PD_SERVICE"),
				Type: "service_reference",
			},
			Body: pagerduty.APIDetails{
				Type:    "incident_body",
				Details: body,
			},
		},
	})
	if err != nil {
		log.Println("Erorr creating incident.  Reason:")
		log.Println(err.Error())
	}
}
