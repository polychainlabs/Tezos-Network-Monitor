package tzrpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// EndorsingRights on the protocol
type EndorsingRights struct {
	data []map[string]interface{}
}

// GetEndorsingRights from the network
// Schema defined here: https://tezos.gitlab.io/alphanet/api/rpc.html#get-block-id-helpers-endorsing-rights
func GetEndorsingRights(level int64) (*EndorsingRights, error) {
	// Get Payload
	resp, err := http.Get(fmt.Sprintf("%v/chains/main/blocks/head/helpers/endorsing_rights?level=%v", os.Getenv("NODE_URL"), level))
	if err != nil {
		log.Println("[Endorsing Rights] Unable to query endpoint")
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[Endorsing Rights] Unable to read response")
		return nil, err
	}

	// Parse Body
	// var parsed []map[string]interface{}
	er := EndorsingRights{}
	err = json.Unmarshal(body, &er.data)
	if err != nil {
		log.Println("[Endorsing Rights] Unable to parse json payload")
		return nil, err
	}

	return &er, nil
}

// Slots that the delegate is allowed to vote
func (er *EndorsingRights) Slots(delegate string) []int64 {
	var slots []int64
	for _, rights := range er.data {
		if rights["delegate"] == delegate {
			for _, slot := range rights["slots"].([]interface{}) {
				slots = append(slots, int64(slot.(float64)))
			}
		}
	}
	return slots
}
