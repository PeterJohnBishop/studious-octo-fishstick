package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func SendRequest(method string, url string, token *string, body *string) ([]byte, error) {

	var requestBody io.Reader
	if body != nil {
		requestBody = strings.NewReader(*body)
	} else {
		requestBody = nil
	}

	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return []byte{}, err
	}
	if token != nil {
		req.Header.Add("Authorization", *token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("API error: %s\nResponse: %s", resp.Status, string(respBody))
	}

	return respBody, nil
}
