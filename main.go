package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mgolfam/gogutils/httpclient"
)

func main() {
	curlCommand := `curl --location --request POST 'https://example.com/v2/auth/token' \
	--header 'Authorization: Basic VmpFNlpucHJWWkRSVTVVUlZJNlJWaFU6YTBkMWREZ3pTbE09' \
	--header 'Content-Type: application/x-www-form-urlencoded' \
	--header 'Cookie: visid_incap_2768614=V3qIDBUB8HdgQx3eqHvW' \
	--data-urlencode 'grant_type=client_credentials'`

	// Call the parser
	config, err := httpclient.ParseCurlCommand(curlCommand)
	if err != nil {
		log.Fatalf("Failed to parse curl command: %v", err)
	}

	// Output the resulting configuration
	fmt.Println("Method:", config.Method)
	fmt.Println("URL:", config.URL)
	fmt.Println("Headers:")
	for key, value := range config.Headers {
		fmt.Printf("  %s: %s\n", key, value)
	}
	fmt.Println("Body:", string(config.Body))
	fmt.Println("Timeout:", config.Timeout)
	fmt.Println("Use Proxy:", config.UseProxy)
	if config.UseProxy {
		fmt.Println("Proxy Address:", config.Headers["Proxy"])
	}

	jj, err := json.Marshal(config)
	fmt.Println(string(jj), err)
}
