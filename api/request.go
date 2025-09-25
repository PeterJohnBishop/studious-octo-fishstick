package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Header struct {
	Key   string
	Value string
}

func SendRequest(method string, headers []Header, url string, body *string) ([]byte, error) {
	var requestBody io.Reader
	if body != nil {
		requestBody = strings.NewReader(*body)
	}

	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for _, h := range headers {
		req.Header.Set(h.Key, h.Value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: %s\nResponse: %s", resp.Status, string(respBody))
	}

	return respBody, nil
}
