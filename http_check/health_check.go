package http_check

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func CheckHealth(url string, method string, expectedStatusCode int, expectedBody string) bool {
	resp, err := http.DefaultClient.Do(&http.Request{
		Method: func() string {
			if method == "" {
				return "GET"
			}
			return strings.ToUpper(method)
		}(),
		URL: parseURL(url), // Convert string to *url.URL
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	// Check for expected status code if provided, otherwise check for http.StatusOK
	if expectedStatusCode != 0 {
		if resp.StatusCode != expectedStatusCode {
			fmt.Printf("Received status code: %d, expected: %d\n", resp.StatusCode, expectedStatusCode)
			return false
		}
	} else if resp.StatusCode != http.StatusOK {
		fmt.Printf("Received status code: %d, expected: %d\n", resp.StatusCode, http.StatusOK)
		return false
	}

	// Check for expected body if provided
	if expectedBody != "" {
		bodyBytes, _ := io.ReadAll(resp.Body) // Read the body
		bodyString := string(bodyBytes)
		if !strings.Contains(bodyString, expectedBody) {
			fmt.Printf("Received body: %s, does not contains: %s\n", bodyString, expectedBody)
			return false
		}
	}

	return true
}

// Function to parse string to *url.URL
func parseURL(rawURL string) *url.URL {
	parsedURL, _ := url.Parse(rawURL) // Handle error as needed
	return parsedURL
}
