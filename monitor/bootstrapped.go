package monitor

import (
	"fmt"

	"github.com/nlopes/slack"
	"gitlab.com/polychainlabs/tezos-network-monitor/alert"
	"gitlab.com/polychainlabs/tezos-network-monitor/tzrpc"
)

// CheckNode health and alert if anything is wrong
func (m *Monitor) CheckNode() {
	bootstrapped, err := tzrpc.GetBootstrapped()
	m.check(err)

	// Slack if lag > 5 minutes
	if bootstrapped.Lag > 60*5 {
		alert.PostSlack(&slack.WebhookMessage{
			Text: fmt.Sprintf("High network lag: `%v minutes`", int(bootstrapped.Lag)/60),
		})
	}
	// Page if lag > 60 minutes
	if bootstrapped.Lag > 60*60 {
		alert.Page(
			fmt.Sprintf("High Tezos Network Lag"),
			fmt.Sprintf("Las is %v minutes.  Has the Tezos network halted or is this node disconnected from the network?", int(bootstrapped.Lag)/60))
	}
}
