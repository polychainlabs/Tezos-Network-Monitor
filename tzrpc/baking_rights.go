package tzrpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// BakingRights on the protocol
type BakingRights struct {
	data []map[string]interface{}
}

// GetBakingRights from the network
// Schema defined here: https://tezos.gitlab.io/alphanet/api/rpc.html#get-block-id-helpers-baking-rights
func GetBakingRights(level int64) (*BakingRights, error) {
	// Get Payload
	resp, err := http.Get(fmt.Sprintf("%v/chains/main/blocks/head/helpers/baking_rights?level=%v", os.Getenv("NODE_URL"), level))
	if err != nil {
		log.Println("[Baking Rights] Unable to query endpoint")
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[Baking Rights] Unable to read response")
		return nil, err
	}

	// Parse Body
	er := BakingRights{}
	err = json.Unmarshal(body, &er.data)
	if err != nil {
		log.Println("[Baking Rights] Unable to parse json payload")
		return nil, err
	}

	return &er, nil
}

// GetBakingPriority of this delegate at this level
func (br *BakingRights) GetBakingPriority(delegate string) int64 {
	priority := int64(-1)
	for _, rights := range br.data {
		if rights["delegate"] == delegate {
			priority = int64(rights["priority"].(float64))
		}
	}
	return priority
}
