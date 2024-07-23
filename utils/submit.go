package utils

import (
	"encoding/json"
	"fmt"

	"github.com/xyield/xrpl-go/client/websocket"
)

// SubmitRequest struct and its methods
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

// New function to submit transactions
func SubmitTransaction(serverURL string, signedTx string) ([]byte, error) {
	ws := websocket.NewWebsocketClient(&websocket.WebsocketConfig{
		URL: serverURL,
	})
	res, err := ws.SendRequest(
		SubmitRequest{
			TxBlob:   signedTx,
			FailHard: false,
		},
	)
	if err != nil {
		return []byte(""), err
	}
	resJSON, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		// Handle error
		fmt.Println("Error marshalling JSON:", err)
		return []byte(""), err
	}

	return resJSON, nil

}
