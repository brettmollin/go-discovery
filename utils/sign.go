package utils

import (
	"encoding/hex"
	"fmt"

	binarycodec "github.com/xyield/xrpl-go/binary-codec"
	"github.com/xyield/xrpl-go/keypairs"
)

/**
 * Sign a transaction
 */
func SignTx(wallet map[string]interface{}, tx_json map[string]any) string {
	// Get the private and public keys from the seed
	privKey, _, _ := keypairs.DeriveKeypair(wallet["seed"].(string), false)

	// Encode the transaction for signing
	encodedTx, _ := binarycodec.EncodeForSigning(tx_json)
	fmt.Println("encodedTx:", encodedTx)

	// Decode the encoded transaction to hexadecimal format
	hexTx, _ := hex.DecodeString(encodedTx)
	fmt.Println("hexTx:", hexTx)

	// Sign the transaction using the private key
	signature, _ := keypairs.Sign(string(hexTx), privKey)
	fmt.Println("signature:", signature)

	// Add the transaction signature to the transaction map
	tx_json["TxnSignature"] = signature

	// Encode the signed transaction
	signedTx, err := binarycodec.Encode(tx_json)
	if err != nil {
		fmt.Println("Error encoding transaction:", err)
	}
	fmt.Println("signedTx:", signedTx)

	return signedTx
}
