package monitor

import (
	"context"
	"log"

	"gitlab.com/polychainlabs/tezos-network-monitor/alert"
	"gitlab.com/polychainlabs/tezos-network-monitor/tzrpc"
)

// Monitor Base
type Monitor struct {
	ctx       context.Context
	addresses []string
	aliases   map[string]string
	whitelist map[string][]string
}

// New monitor
func New(ctx context.Context, addresses []string, aliases map[string]string, whitelist map[string][]string) *Monitor {
	m := Monitor{
		ctx:       ctx,
		addresses: addresses,
		aliases:   aliases,
		whitelist: whitelist,
	}

	return &m
}

// helper to get the latest block
func (m *Monitor) getCurrentBlock() *tzrpc.Block {
	bootstrapped, err := tzrpc.GetBootstrapped()
	m.check(err)

	// Get latest Block
	block, err := tzrpc.GetBlock(bootstrapped.Block, 0)
	m.check(err)
	return block
}

func (m *Monitor) alias(address string) string {
	return alert.Alias(m.aliases, address)
}

func (m *Monitor) logError(err error) {
	log.Println(err)
}

func (m *Monitor) check(err error) {
	m.fatalError(err)
}

func (m *Monitor) fatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
