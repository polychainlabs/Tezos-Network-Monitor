package monitor

import (
	"fmt"
	"log"

	"github.com/nlopes/slack"
	"gitlab.com/polychainlabs/tezos-network-monitor/alert"
	"gitlab.com/polychainlabs/tezos-network-monitor/storage"
	"gitlab.com/polychainlabs/tezos-network-monitor/tzrpc"
)

// CheckBlocks and alerts if any error conditions are met
func (m *Monitor) CheckBlocks() {
	currentBlock := m.getCurrentBlock()

	// Get last checked level from firestore
	lastRecordedLevel := storage.GetLastRecordedBlockLevel()

	// Ignore blocks more than 1 cycle ago
	if lastRecordedLevel == -1 || lastRecordedLevel < currentBlock.Level()-4096 {
		minLevel := currentBlock.Level() - 2
		m.logError(fmt.Errorf("[Block]\tLast recorded level of %v is too low.  Resetting to %v", lastRecordedLevel, minLevel))
		lastRecordedLevel = minLevel
	}

	// Analyze all new blocks
	for level := lastRecordedLevel + 1; level < currentBlock.Level(); level++ {
		debugln("[Block]\tAnalyzing level ", level)

		// Get block at level
		block, err := tzrpc.GetBlock(currentBlock.Hash(), currentBlock.Level()-level)
		m.check(err)

		// Alert when transactions are sent or received
		for _, tx := range block.Transactions() {
			for _, address := range m.addresses {
				if address == tx.Source {
					// Slack when transactions sent _from_ your address
					alert.PostSlack(&slack.WebhookMessage{
						Text: fmt.Sprintf("*Sent* `%v`ꜩ from `%v` to `%v` with `%v`ꜩ fee",
							tx.Amount/1e6, m.alias(tx.Source), m.alias(tx.Destination), tx.Fee/1e6),
					})
					// Page if destination address is not whitelisted
					if !m.isDestinationWhitelisted(tx.Source, tx.Destination) {
						title := fmt.Sprintf("Sent %vꜩ from %v", tx.Amount/1e6, m.alias(tx.Source))
						body := fmt.Sprintf("To %v with fee %vꜩ at level %v",
							m.alias(tx.Destination), tx.Fee/1e6, level)
						alert.Page(title, body)
					}
				}
				if address == tx.Destination {
					// Slack when transactions sent _to_ your address
					alert.PostSlack(&slack.WebhookMessage{
						Text: fmt.Sprintf("*Received* `%v`ꜩ at `%v`", tx.Amount/1e6, m.alias(tx.Destination)),
					})
				}
			}
		}

		// Slack when delegations received
		for _, delegation := range block.Delegations() {
			for _, address := range m.addresses {
				if address == delegation.Delegate {
					amount := getBalanceString(delegation.Source)
					alert.PostSlack(&slack.WebhookMessage{
						Text: fmt.Sprintf("*Delegation* `%v` delegated `%vꜩ` to `%v`",
							m.alias(delegation.Source), amount, m.alias(delegation.Delegate)),
					})
				}
				if address == delegation.Source {
					amount := getBalanceString(delegation.Source)
					alert.PostSlack(&slack.WebhookMessage{
						Text: fmt.Sprintf("*Delegation* We delegated `%vꜩ`from `%v` to `%v`",
							amount, m.alias(delegation.Source), m.alias(delegation.Delegate)),
					})
				}
			}
		}

		// Slack when delegations received through originations
		for _, origination := range block.Originations() {
			for _, address := range m.addresses {
				if address == origination.Delegate {
					alert.PostSlack(&slack.WebhookMessage{
						Text: fmt.Sprintf("*Delegation* `%v` delegated `%vꜩ` to `%v` through an origination",
							origination.Source, origination.Balance.String(), m.alias(origination.Delegate)),
					})
				}
				if address == origination.Source {
					alert.PostSlack(&slack.WebhookMessage{
						Text: fmt.Sprintf("*Origination* We originated a contract from `%v`",
							m.alias(origination.Source)),
					})
				}
			}
		}

		// Alert on double baking
		for _, double := range block.DoubleBakings() {
			// Slack if anyone has double baked
			alert.PostSlack(&slack.WebhookMessage{
				Text: fmt.Sprintf("*Double Baking* found at level `%v`. `%vꜩ` slashed", double.Level, double.SlashedAmount),
			})
			for _, address := range m.addresses {
				if address == double.SlashedBaker {
					// Page if you've double baked :(
					alert.PostSlack(&slack.WebhookMessage{
						Text: fmt.Sprintf("*WE DOUBLE BAKED* At level `%v` by baker `%v`. `%vꜩ` slashed. SHUT THIS BAKER DOWN NOW.", double.Level, double.SlashedBaker, double.SlashedAmount),
					})

					title := fmt.Sprintf("WE DOUBLE BAKED with %v", address)
					body := fmt.Sprintf("%v was just slashed at level %v.  SHUT THIS BAKER DOWN NOW, AND STAY OFFLINE FOR THE REMAINED OF THE CYCLE.", double.SlashedAmount, level)
					alert.Page(title, body)
				}
			}
		}

		// Alert on double endorsements
		for _, double := range block.DoubleEndorsements() {
			// Slack if anyone has double endorsed
			alert.PostSlack(&slack.WebhookMessage{
				Text: fmt.Sprintf("*Double Endorsement* found at level `%v`. `%vꜩ` slashed", double.Level, double.SlashedAmount),
			})
			for _, address := range m.addresses {
				if address == double.SlashedBaker {
					// Page if you've double endorsed :(
					alert.PostSlack(&slack.WebhookMessage{
						Text: fmt.Sprintf("*WE DOUBLE ENDORSED* At level `%v` by endorser `%v`. `%vꜩ` slashed. SHUT THIS ENDORSER DOWN NOW", double.Level, double.SlashedBaker, double.SlashedAmount),
					})

					title := fmt.Sprintf("WE DOUBLE ENDORSED with %v at level %v", address, level)
					body := fmt.Sprintf("%v was just slashed.  SHUT THIS ENDORSER DOWN NOW, AND STAY OFFLINE FOR THE REMAINED OF THE CYCLE", double.SlashedAmount)
					alert.Page(title, body)
				}
			}
		}

		// Save
		storage.RecordBlock(level, block.Hash())
	}
}

func (m *Monitor) isDestinationWhitelisted(source string, destination string) bool {
	if _, ok := m.whitelist[source]; !ok {
		return false
	}
	for _, address := range m.whitelist[source] {
		if address == destination {
			return true
		}
	}
	return false
}

func getBalanceString(pkh string) string {
	balance, err := tzrpc.GetBalance(pkh)
	if err != nil {
		log.Println(err)
		return ""
	}
	return balance.String()
}
