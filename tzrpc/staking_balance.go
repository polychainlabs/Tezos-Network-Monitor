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

// GetStakingBalance of this public key hash.  Returned
// amount is in *full* Tez. Schema defined here:
// https://tezos.gitlab.io/alphanet/api/rpc.html#get-block-id-context-delegates-pkh-staking-balance
func GetStakingBalance(pkh string) (*big.Int, error) {
	// Get Payload
	resp, err := http.Get(fmt.Sprintf("%v/chains/main/blocks/head/context/delegates/%v/staking_balance", os.Getenv("NODE_URL"), pkh))
	if err != nil {
		log.Println("[Staking Balance] Unable to query endpoint")
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[Staking Balance] Unable to read response")
		return nil, err
	}

	// Parse Body
	balanceString := strings.Trim(string(body), "\"\n")
	balance, _ := new(big.Int).SetString(balanceString, 10)
	if balance == nil {
		return nil, errors.New("[Staking Balance] Could not parse balance int")
	}
	// Convert from uTez to Tez
	balance.Div(balance, new(big.Int).SetInt64(1000000))
	return balance, nil
}
