package main

import (
	"context"
	"log"
	"time"

	"gitlab.com/polychainlabs/tezos-network-monitor/monitor"
)

func main() {
	log.Println("Starting up...")

	// Setup
	ctx := context.Background()
	c := loadConfig("./config.yaml")
	addresses := append(c.Bakers, c.Delegators...)

	// Monitor
	monitor := monitor.New(ctx, addresses, c.Aliases, c.Whitelist)

	for {
		// Alert if network has stopped
		monitor.CheckNode()

		// Monitor Blocks
		monitor.CheckBlocks()

		// Monitor Endoring and Baking Trends
		for _, d := range c.Bakers {
			monitor.CheckBaking(d)
			monitor.CheckEndorsing(d)
		}

		// Sleep
		sleepSeconds := 5 * time.Second
		log.Printf("Successfully caught up.  Sleeping for %v\n", sleepSeconds)
		time.Sleep(sleepSeconds)
	}
}
