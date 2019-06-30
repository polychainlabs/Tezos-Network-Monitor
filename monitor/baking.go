package monitor

import (
	"fmt"
	"log"

	"github.com/nlopes/slack"
	"gitlab.com/polychainlabs/tezos-network-monitor/alert"
	"gitlab.com/polychainlabs/tezos-network-monitor/storage"
	"gitlab.com/polychainlabs/tezos-network-monitor/tzrpc"
)

// CheckBaking performance for this delegate
func (m *Monitor) CheckBaking(delegate string) {
	currentBlock := m.getCurrentBlock()

	// Get last checked level from firestore
	lastRecordedLevel := storage.GetLastRecordedBakeLevel(delegate)

	// Baking rights are only availabe so many levels behind, so start no earlier than 1 cycle ago
	if lastRecordedLevel == -1 || lastRecordedLevel < currentBlock.Level()-4096 {
		minLevel := currentBlock.Level() - 2
		m.logError(fmt.Errorf("[Baking]\tLast recorded level of %v is too low.  Resetting to %v",
			lastRecordedLevel, minLevel))
		lastRecordedLevel = minLevel
	}

	// Analyze baking rights
	for level := lastRecordedLevel + 1; level < currentBlock.Level(); level++ {
		debugf("[Baking]\tAnalyzing level %v for %v\n", level, delegate)
		// Get rights at level
		bakingRights, err := tzrpc.GetBakingRights(level)
		m.check(err)
		delegateRights := bakingRights.GetBakingPriority(delegate)

		// Get remainder of metadata from blocks
		block, err := tzrpc.GetBlock(currentBlock.Hash(), currentBlock.Level()-level)
		m.check(err)

		blockHash := block.Hash()
		bakerPriority := block.BakerPriority()

		// Alert on misses
		if delegateRights >= 0 && bakerPriority > delegateRights {
			// Slack when you miss a block
			alert.PostSlack(&slack.WebhookMessage{
				Text: fmt.Sprintf("*Missed Block* at level `%v` by `%v`", level, m.alias(delegate)),
			})

			// Machine parseable logline
			log.Printf("baker=%v level=%v miss=1\n", delegate, level)

			// Page if we've missed a lot this cycle
			m.checkBakingTrends(delegate)
		} else if bakerPriority == delegateRights {
			// Machine parseable logline when we've baked
			log.Printf("baker=%v level=%v miss=0\n", delegate, level)
		}

		// Save to Datastore
		storage.RecordBaking(delegate, level, delegateRights, bakerPriority, blockHash)
	}
}

// checkBakingTrends and page if you've missed a lot this cycle
func (m *Monitor) checkBakingTrends(delegate string) {
	// Last X Misses
	misses, _ := storage.GetCycleBakeMissCount(delegate)

	// Page if miss > 2 bakings per cycle
	if misses > 2 {
		title := fmt.Sprintf("Missed many blocks by %v", delegate)
		body := fmt.Sprintf("Missed %v blocks.  Is this baker online?", misses)
		alert.Page(title, body)
	}
}
