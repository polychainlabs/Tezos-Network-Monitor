package tzrpc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
)

// GetBalance of this public key hash. Returned
// amount is in *full* Tez. Schema defined here:
// https://tezos.gitlab.io/alphanet/api/rpc.html#get-block-id-context-contracts-contract-id-balance
func GetBalance(pkh string) (*big.Int, error) {
	// Get Payload
	var path string
	if len(pkh) > 2 && strings.HasPrefix(pkh, "KT1") {
		path = "chains/main/blocks/head/context/contracts"
	} else if len(pkh) > 1 && strings.HasPrefix(pkh, "tz") {
		path = "chains/main/blocks/head/context/delegates"
	} else {
		return nil, errors.New("Invalid pkh format")
	}
	resp, err := http.Get(fmt.Sprintf("%v/%v/%v/balance", os.Getenv("NODE_URL"), path, pkh))
	if err != nil {
		log.Println("[Balance] Unable to query endpoint")
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[Balance] Unable to read response")
		return nil, err
	}

	// Parse Body
	balanceString := strings.Trim(string(body), "\"\n")
	balance, _ := new(big.Int).SetString(balanceString, 10)
	if balance == nil {
		return nil, errors.New("[Balance] Could not parse balance int")
	}
	// Convert from uTez to Tez
	balance.Div(balance, new(big.Int).SetInt64(1000000))
	return balance, nil
}
