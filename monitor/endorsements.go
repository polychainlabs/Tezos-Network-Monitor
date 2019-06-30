package monitor

import (
	"fmt"
	"log"

	"github.com/nlopes/slack"
	"gitlab.com/polychainlabs/tezos-network-monitor/alert"
	"gitlab.com/polychainlabs/tezos-network-monitor/storage"
	"gitlab.com/polychainlabs/tezos-network-monitor/tzrpc"
)

// CheckEndorsing performance for this delegate
func (m *Monitor) CheckEndorsing(delegate string) {
	currentBlock := m.getCurrentBlock()

	// Get last checked level from firestore
	lastRecordedLevel := storage.GetLastRecordedEndorsementLevel(delegate)

	// Endorsing rights are only availabe so many levels behind, so start no earlier than 1 cycle ago
	if lastRecordedLevel == -1 || lastRecordedLevel < currentBlock.Level()-4096 {
		minLevel := currentBlock.Level() - 2
		m.logError(fmt.Errorf("[Endorsing]\tLast recorded level of %v is too low.  Resetting to %v",
			lastRecordedLevel, minLevel))
		lastRecordedLevel = minLevel
	}

	// Analyze endoring rights
	for level := lastRecordedLevel + 1; level < currentBlock.Level(); level++ {
		debugf("[Endorsing]\tAnalyzing level %v for %v\n", level, delegate)
		// Get rights at level
		endorsingRights, err := tzrpc.GetEndorsingRights(level)
		m.check(err)
		rights := endorsingRights.Slots(delegate)

		// Get endorsements from block `n-1`
		block, err := tzrpc.GetBlock(currentBlock.Hash(), currentBlock.Level()-level-1)
		m.check(err)
		endorsements, err := block.Endorsements(delegate)
		m.check(err)

		// Get block hash from block `n`
		block, err = tzrpc.GetBlock(currentBlock.Hash(), currentBlock.Level()-level)
		m.check(err)
		hash := block.Hash()

		// Alert on misses
		if len(rights) > len(endorsements) {
			alert.PostSlack(&slack.WebhookMessage{
				Text: fmt.Sprintf("*Missed Endorsement* at level `%v` with baker `%v`",
					level, m.alias(delegate)),
			})

			// Page if we've missed a lot this cycle
			m.checkEndorsingTrends(delegate)
		}

		// Machine parseable logline
		log.Printf("endorser=%v endorsements=%v misses=%v level=%v\n",
			delegate, len(endorsements), len(rights)-len(endorsements), level)

		// Save to Datastore
		storage.RecordEndorsement(delegate, level, rights, endorsements, hash)
	}
}

// checkEndorsingTrends and alert if we're missing a lot
func (m *Monitor) checkEndorsingTrends(delegate string) {
	previousLevels := 20
	endorsements := storage.GetEndorsements(delegate, previousLevels)

	nMisses := 0
	var level int64
	for _, e := range endorsements {
		if e.Misses > 0 {
			nMisses++
		}
		level = e.Level
	}

	// Page if miss 2 or more of last `previousLevels` endorsements
	if nMisses >= 2 {
		title := fmt.Sprintf("Missed %v of last %v endorsements for %v",
			nMisses, previousLevels, delegate)
		body := fmt.Sprintf("Is this baker online? From level %v", level)
		alert.Page(title, body)
	}

	// Cycle Misses
	cycleMisses := storage.GetCycleEndorsementMissCount(delegate)
	// Page if miss more than 5 per cycle
	if cycleMisses > 5 {
		title := fmt.Sprintf("Missed many endorsements by %v", delegate)
		body := fmt.Sprintf("Missed %v endorsements.  Is this baker online?", cycleMisses)
		alert.Page(title, body)
	}
}
