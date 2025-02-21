package main

import (
	"io"
	"net/http"
)

func Any(arr []bool) bool {
	for _, v := range arr {
		if v {
			return true
		}
	}
	return false
}

func getPublicIP() (string, error) {
	// Use an external service to fetch the public IP
	resp, err := http.Get("https://api64.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Convert response to string
	return string(body), nil
}
