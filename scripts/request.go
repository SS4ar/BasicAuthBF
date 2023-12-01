package lib

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// SendHTTPRequest sends a GET request to the server with an encoded token value, ignoring SSL errors
func SendHTTPRequest(url string, encodedToken string, wg *sync.WaitGroup) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	// Set the Authorization header with the encoded token
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedToken))

	// Create a custom HTTP client with InsecureSkipVerify set to true
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: tr,
	}

	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}

	// Close the response body
	resp.Body.Close()

	// Decrement the wait group counter when the function completes
	defer wg.Done()

	// Check the response status code
	if resp.StatusCode == 200 {
		// If the status code is 200, print success message and exit
		fmt.Println()
		PrintSuccess("Found valid credentials", Base64Decode(encodedToken), false, false)
		fmt.Println()
		fmt.Println()
		os.Exit(0)
	} else {
		// If the status code is not 200, print failure message
		PrintFailed("failed", Base64Decode(encodedToken), false, true)
	}
}
