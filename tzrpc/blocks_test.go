package tzrpc

import (
	"io/ioutil"
	"log"
	"math/big"
	"testing"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getBlock(path string) *Block {
	body, err := ioutil.ReadFile(path)
	check(err)
	block, err := parseBlock(body)
	check(err)
	return block
}

func TestDoubleBaking(t *testing.T) {
	block := getBlock("../tests/double_baking.json")

	doubles := block.DoubleBakings()

	if len(doubles) != 1 {
		log.Println("Failure: Expected 1 double bakings but found", len(doubles))
		t.Fail()
	}

	if doubles[0].SlashedBaker != "tz1Kt4P8BCaP93AEV4eA7gmpRryWt5hznjCP" {
		log.Println("Failure: Incorrect slashed baker.  Received", doubles[0].SlashedBaker)
		t.Fail()
	}
	if doubles[0].RewardedBaker != "tz1WnfXMPaNTBmH7DBPwqCWs9cPDJdkGBTZ8" {
		log.Println("Failure: Incorrect rewarded baker.  Received", doubles[0].RewardedBaker)
		t.Fail()
	}
	if doubles[0].SlashedAmount != -64000000 {
		log.Println("Failure: Incorrect slashed amount.  Received", doubles[0].SlashedAmount)
		t.Fail()
	}
	if doubles[0].Level != 74 {
		log.Println("Failure: Incorrect slashed level.  Received", doubles[0].Level)
		t.Fail()
	}
}

func TestDelegation(t *testing.T) {
	block := getBlock("../tests/delegation.json")

	delegations := block.Delegations()

	if len(delegations) != 1 {
		log.Println("Failure: Expected 1 delegation but found", len(delegations))
		t.Fail()
	}

	if delegations[0].Source != "KT1D58hJ8msoXAsjmCYGjkCFRB9APx5VvrnL" {
		log.Println("Failure: Incorrect delegation source.  Received", delegations[0].Source)
		t.Fail()
	}
	if delegations[0].Delegate != "tz2PdGc7U5tiyqPgTSgqCDct94qd6ovQwP6u" {
		log.Println("Failure: Incorrect Delegate.  Received", delegations[0].Delegate)
		t.Fail()
	}

}

func TestOrigination(t *testing.T) {
	block := getBlock("../tests/origination.json")

	originations := block.Originations()

	if len(originations) != 1 {
		log.Println("Failure: Expected 1 delegation but found", len(originations))
		t.Fail()
	}

	if originations[0].Source != "tz1X8KDWYVH4rhET5k1aJuU9q6h2vR5kjezG" {
		log.Println("Failure: Incorrect delegation source.  Received", originations[0].Source)
		t.Fail()
	}
	if originations[0].Delegate != "tz2PdGc7U5tiyqPgTSgqCDct94qd6ovQwP6u" {
		log.Println("Failure: Incorrect Delegate.  Received", originations[0].Delegate)
		t.Fail()
	}
	if originations[0].Balance.Cmp(new(big.Int).SetInt64(1050)) != 0 {
		log.Println("Failure: Incorrect Amount.  Received", originations[0].Balance)
		t.Fail()
	}

}
