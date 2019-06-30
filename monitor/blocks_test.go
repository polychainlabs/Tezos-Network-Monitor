package monitor

import (
	"testing"
)

func TestWhitelist(t *testing.T) {

	whitelist := map[string][]string{
		"tz1a": []string{
			"tz1b",
			"KT1c",
		},
	}

	m := Monitor{
		whitelist: whitelist,
	}

	// Should Succeed
	if !m.isDestinationWhitelisted("tz1a", "tz1b") {
		t.Fail()
	}
	// Should Fail
	if m.isDestinationWhitelisted("tz1a", "tz1a") {
		t.Fail()
	}
	if m.isDestinationWhitelisted("tz1d", "tz1d") {
		t.Fail()
	}
}
