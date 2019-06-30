package tzrpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Bootstrapped RPC response
type Bootstrapped struct {
	Block     string
	Timestamp time.Time
	Lag       float64
}

// GetBootstrapped struct.  Raw response is like:
//  {
//  	"block": "BLmyFBmkuqUvPpEb3GUzX2snBTqcGMD4UyKqTF1yNTGA5mRBi6E",
//  	"timestamp": "2019-03-20T06:31:36Z"
//  }
func GetBootstrapped() (*Bootstrapped, error) {
	// Get Payload with 15s timeout
	client := http.Client{
		Timeout: time.Duration(15 * time.Second),
	}

	resp, err := client.Get(fmt.Sprintf("%v/monitor/bootstrapped", os.Getenv("NODE_URL")))
	if err != nil {
		log.Println("[Bootstrapped] Unable to query endpoint: ", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[Bootstrapped] Unable to read response: ", err)
		return nil, err
	}

	// Parse Body
	var bootstrapped map[string]string
	err = json.Unmarshal(body, &bootstrapped)
	if err != nil {
		log.Println("[Bootstrapped] Unable to parse json payload: ", err)
		return nil, err
	}
	timestamp, err := time.Parse(time.RFC3339, bootstrapped["timestamp"])
	if err != nil {
		log.Println("[Bootstrapped] Unable to parse timestamp: ", err)
		return nil, err
	}
	secondsBehind := time.Now().Sub(timestamp).Seconds()

	return &Bootstrapped{
		Block:     bootstrapped["block"],
		Timestamp: timestamp,
		Lag:       secondsBehind,
	}, nil
}
