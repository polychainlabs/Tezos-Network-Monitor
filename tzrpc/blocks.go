package tzrpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
)

// Block data
type Block struct {
	data map[string]interface{}
}

// GetBlock data
func GetBlock(blockHash string, offset int64) (*Block, error) {
	// Get Payload
	resp, err := http.Get(fmt.Sprintf("%v/chains/main/blocks/%v~%v", os.Getenv("NODE_URL"), blockHash, offset))
	if err != nil {
		log.Println("[Block Endorsements] Unable to query endpoint")
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[Block Endorsements] Unable to read response")
		return nil, err
	}

	return parseBlock(body)
}

func parseBlock(body []byte) (*Block, error) {
	// Parse Body
	var block Block
	err := json.Unmarshal(body, &block.data)

	if err != nil {
		log.Println("[Block Endorsements] Unable to parse json payload", err)
		return nil, err
	}

	return &block, nil
}

// Level of the current block
func (block *Block) Level() int64 {
	header := block.data["header"].(map[string]interface{})
	return int64(header["level"].(float64))
}

// Hash of the current block
func (block *Block) Hash() string {
	return block.data["hash"].(string)
}

// Baker of the current block
func (block *Block) Baker() string {
	metadata := block.data["metadata"].(map[string]interface{})
	return metadata["baker"].(string)
}

// BakerPriority that baked this block
func (block *Block) BakerPriority() int64 {
	header := block.data["header"].(map[string]interface{})
	return int64(header["priority"].(float64))
}

// Endorsements for the current block
func (block *Block) Endorsements(delegate string) ([]int64, error) {
	operations := block.data["operations"].([]interface{})

	var slots []int64
	for _, outer := range operations {
		for _, inner := range outer.([]interface{}) {
			operation := inner.(map[string]interface{})
			contents := operation["contents"].([]interface{})
			for _, c := range contents {
				typedContents := c.(map[string]interface{})
				if typedContents["kind"] != "endorsement" {
					continue
				}
				metadata := typedContents["metadata"].(map[string]interface{})
				if metadata["delegate"] != delegate {
					continue
				}
				for _, slotID := range metadata["slots"].([]interface{}) {
					slots = append(slots, int64(slotID.(float64)))
				}
			}
		}
	}
	return slots, nil
}

// Transaction data
type Transaction struct {
	Source      string
	Destination string
	Amount      int
	Fee         int
}

// Transactions in the block
func (block *Block) Transactions() []Transaction {
	tx := []Transaction{}

	operations := block.data["operations"].([]interface{})
	for _, outer := range operations {
		for _, inner := range outer.([]interface{}) {
			operation := inner.(map[string]interface{})
			contents := operation["contents"].([]interface{})
			for _, c := range contents {
				typedContents := c.(map[string]interface{})
				if typedContents["kind"] == "transaction" {
					amount, _ := strconv.Atoi(typedContents["amount"].(string))
					fee, _ := strconv.Atoi(typedContents["fee"].(string))

					tx = append(tx, Transaction{
						Source:      typedContents["source"].(string),
						Destination: typedContents["destination"].(string),
						Amount:      amount,
						Fee:         fee,
					})
				}
			}
		}
	}
	return tx
}

// DoubleBaking data
type DoubleBaking struct {
	SlashedBaker  string
	RewardedBaker string
	SlashedAmount int
	Level         int
}

// DoubleBakings in the block
func (block *Block) DoubleBakings() []DoubleBaking {
	doubles := []DoubleBaking{}

	operations := block.data["operations"].([]interface{})
	for _, outer := range operations {
		for _, inner := range outer.([]interface{}) {
			operation := inner.(map[string]interface{})
			contents := operation["contents"].([]interface{})
			for _, c := range contents {
				typedContents := c.(map[string]interface{})
				if typedContents["kind"] == "double_baking_evidence" {
					typedMetadata := typedContents["metadata"].(map[string]interface{})
					balanceUpdates := typedMetadata["balance_updates"].([]interface{})

					double := DoubleBaking{}

					for _, update := range balanceUpdates {
						typedUpdate := update.(map[string]interface{})
						if typedUpdate["category"] == "deposits" {
							double.SlashedBaker = typedUpdate["delegate"].(string)
							amount, _ := strconv.Atoi(typedUpdate["change"].(string))
							double.SlashedAmount = amount
							double.Level = int(typedUpdate["level"].(float64))
						} else if typedUpdate["category"] == "rewards" {
							double.RewardedBaker = typedUpdate["delegate"].(string)
						}
					}

					doubles = append(doubles, double)
				}
			}
		}
	}
	return doubles
}

// DoubleEndorsement data
type DoubleEndorsement struct {
	SlashedEndorser string
	RewardedBaker   string
	SlashedAmount   int
	Cycle           int
}

// DoubleEndorsements in the block
func (block *Block) DoubleEndorsements() []DoubleEndorsement {
	doubles := []DoubleEndorsement{}

	operations := block.data["operations"].([]interface{})
	for _, outer := range operations {
		for _, inner := range outer.([]interface{}) {
			operation := inner.(map[string]interface{})
			contents := operation["contents"].([]interface{})
			for _, c := range contents {
				typedContents := c.(map[string]interface{})
				if typedContents["kind"] == "double_endorsement_evidence" {
					typedMetadata := typedContents["metadata"].(map[string]interface{})
					balanceUpdates := typedMetadata["balance_updates"].([]interface{})

					double := DoubleEndorsement{}

					for _, update := range balanceUpdates {
						typedUpdate := update.(map[string]interface{})
						switch typedUpdate["category"] {
						case "deposits", "fees", "rewards":
							amount, _ := strconv.Atoi(typedUpdate["change"].(string))
							if amount < 0 {
								double.SlashedEndorser = typedUpdate["delegate"].(string)
								double.Cycle = int(typedUpdate["cycle"].(float64))
								amount, _ := strconv.Atoi(typedUpdate["change"].(string))
								double.SlashedAmount -= amount
							} else if amount > 0 {
								double.RewardedBaker = typedUpdate["delegate"].(string)
							}
						}
					}

					doubles = append(doubles, double)
				}
			}
		}
	}
	return doubles
}

// Delegation to a baker
type Delegation struct {
	Source   string
	Delegate string
}

// Delegations in the block
func (block *Block) Delegations() []Delegation {
	delegations := []Delegation{}

	operations := block.data["operations"].([]interface{})
	for _, outer := range operations {
		for _, inner := range outer.([]interface{}) {
			operation := inner.(map[string]interface{})
			contents := operation["contents"].([]interface{})
			for _, c := range contents {
				typedContents := c.(map[string]interface{})
				if typedContents["kind"] == "delegation" {
					// Create a new delegation
					newDelegation := Delegation{}
					// Fields can be missing, so chech to see if they exist
					if typedContents["delegate"] != nil {
						newDelegation.Delegate = typedContents["delegate"].(string)
					}
					if typedContents["source"] != nil {
						newDelegation.Source = typedContents["source"].(string)
					}
					delegations = append(delegations, newDelegation)
				}
			}
		}
	}
	return delegations
}

// Origination to a baker
type Origination struct {
	Source   string
	Delegate string
	Balance  *big.Int
}

// Originations in the block
func (block *Block) Originations() []Origination {
	originations := []Origination{}

	operations := block.data["operations"].([]interface{})
	for _, outer := range operations {
		for _, inner := range outer.([]interface{}) {
			operation := inner.(map[string]interface{})
			contents := operation["contents"].([]interface{})
			for _, c := range contents {
				typedContents := c.(map[string]interface{})
				if typedContents["kind"] == "origination" {
					// Create a new delegation
					newOrigination := Origination{}
					// Fields can be missing, so chech to see if they exist
					if typedContents["delegate"] != nil {
						newOrigination.Delegate = typedContents["delegate"].(string)
					}
					if typedContents["source"] != nil {
						newOrigination.Source = typedContents["source"].(string)
					}

					newOrigination.Balance = &big.Int{}
					if typedContents["balance"] != nil {
						// Parse Body
						balanceString := typedContents["balance"].(string)
						balance, _ := new(big.Int).SetString(balanceString, 10)
						if balance != nil {
							newOrigination.Balance = balance.Div(balance, new(big.Int).SetInt64(1000000))
						}
					}
					originations = append(originations, newOrigination)
				}
			}
		}
	}
	return originations
}
