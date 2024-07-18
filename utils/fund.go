package utils

// Import the required packages
import (
	"fmt"
	//for http requests to faucet
	"encoding/json"
	"io"
	"net/http"
)

func FundWallet() map[string]interface{} {
	const FaucetURL = "https://faucet.altnet.rippletest.net/accounts"

	request, err := http.NewRequest("POST", FaucetURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)

	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)

	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)

	}

	// Parse the JSON into a generic interface{}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error parsing response body:", err)

	}

	return result
}
