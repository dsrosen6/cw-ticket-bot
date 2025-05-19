package webex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseUrl = "https://webexapis.com/v1"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(httpClient *http.Client, apiKey string) *Client {
	return &Client{
		httpClient: httpClient,
		apiKey:     apiKey,
	}
}

func (c *Client) request(ctx context.Context, method, endpoint string, payload io.Reader, target interface{}) error {
	url := fmt.Sprintf("%s/%s", baseUrl, endpoint)
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		return fmt.Errorf("creating the request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending the request: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("reading error response body: %w", err)
		}
		return fmt.Errorf("non-success response code: %d: %s", res.StatusCode, string(data))
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading the response body: %w", err)
	}

	if target != nil {
		if err := json.Unmarshal(data, target); err != nil {
			return fmt.Errorf("unmarshaling the response to json: %w", err)
		}
	}

	return nil
}
