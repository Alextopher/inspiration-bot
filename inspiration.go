package main

import (
	"io"
	"net/http"
)

// Get a link to an inspirational image
func getLink() (string, error) {
	resp, err := http.Get("https://inspirobot.me/api?generate=true")
	if err != nil {
		return "", err
	}

	// get the link out of the body
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	resp.Body.Close()

	return string(bytes), nil
}
