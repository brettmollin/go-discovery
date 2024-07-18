package main

import (
	"encoding/hex"
	"fmt"
	"go-discovery/first-app/utils"

	binarycodec "github.com/xyield/xrpl-go/binary-codec"
	"github.com/xyield/xrpl-go/client/websocket"
	"github.com/xyield/xrpl-go/keypairs"
	"github.com/xyield/xrpl-go/model/client/account"
	"github.com/xyield/xrpl-go/model/transactions/types"
)

// Specify the node to connect to
// devnet wss://s.devnet.rippletest.net:51233/
// testnet wss://s.altnet.rippletest.net:51233
// mainnet wss://s1.ripple.com:443
// to not be rate limited connect to your own node and list IP ad admin in rippled.cfg
const ServerUrl = "wss://s.altnet.rippletest.net:51233/"

type SubmitRequest struct {
	TxBlob   string `json:"tx_blob"`
	FailHard bool   `json:"fail_hard,omitempty"`
}

func (SubmitRequest) Validate() error {
	return nil
}
func (SubmitRequest) Method() string {
	return "submit"
}

/**
 * Example to send a payment
 */
func main() {
	// Create a funded account with test/dev net faucet and return the wallet details in json
	sender_account := utils.FundWallet()
	// Create a funder account for the destination account
	destination_account := utils.FundWallet()

	// Get the address of the funded account
	address := sender_account["account"].(map[string]interface{})["address"].(string)

	//Connect to the ledger to get the next account sequence number to use
	client := websocket.NewClient(&websocket.WebsocketConfig{URL: ServerUrl})
	acr, _, err := client.Account.GetAccountInfo(&account.AccountInfoRequest{Account: types.Address(address)})
	if err != nil {
		panic(err)
	}
	// Print the seeds for the sender and destination accounts
	fmt.Println("private key for sender", sender_account["seed"].(string))
	fmt.Println("private key for destination", destination_account["seed"].(string))

	// Get the private and public keys from the seed
	privKey, pubKey, _ := keypairs.DeriveKeypair(sender_account["seed"].(string), false)

	// Get the distintion address
	DestinationAddress := destination_account["account"].(map[string]interface{})["address"].(string)

	tx := map[string]any{
		"Account":         address,
		"TransactionType": "Payment",
		"Amount":          "20",
		"Destination":     DestinationAddress,
		"Flags":           0,
		"Fee":             "12", //TODO: get live fee from the ledger
		"Sequence":        int(acr.AccountData.Sequence),
		"SigningPubKey":   pubKey,
	}
	// Create a new websocket client to connect to the server
	ws := websocket.NewWebsocketClient(&websocket.WebsocketConfig{
		URL: ServerUrl,
	})

	// Encode the transaction for signing
	encodedTx, _ := binarycodec.EncodeForSigning(tx)
	fmt.Println("encodedTx:", encodedTx)

	// Decode the encoded transaction to hexadecimal format
	hexTx, _ := hex.DecodeString(encodedTx)
	fmt.Println("hexTx:", hexTx)

	// Sign the transaction using the private key
	signature, _ := keypairs.Sign(string(hexTx), privKey)
	fmt.Println("signature:", signature)

	// Add the transaction signature to the transaction map
	tx["TxnSignature"] = signature

	// Encode the signed transaction
	signedTx, err := binarycodec.Encode(tx)
	if err != nil {
		return
	}
	fmt.Println("signedTx:", signedTx)

	// Send the signed transaction to the server for submission
	res, _ := ws.SendRequest(
		SubmitRequest{
			signedTx,
			false,
		})

	fmt.Println("Tx Result", res)
}
