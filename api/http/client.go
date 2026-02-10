package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client interface {
	Get(ctx context.Context, url string, result interface{}) error
	Post(ctx context.Context, url string, data interface{}, result interface{}) error
}

type HTTPClient struct {
	client  *http.Client
	baseURL string
	authKey string
}

func NewHTTPClient(baseURL string, authKey string) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
		authKey: authKey,
	}
}

func (c *HTTPClient) Get(ctx context.Context, path string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
	if err != nil {
		return err
	}

	if c.authKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.authKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *HTTPClient) Post(ctx context.Context, path string, data interface{}, result interface{}) error {
	var body bytes.Buffer
	if data != nil {
		if err := json.NewEncoder(&body).Encode(data); err != nil {
			return fmt.Errorf("failed to encode request: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, &body)
	if err != nil {
		return err
	}

	if c.authKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.authKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
