package connectwise

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseUrl = "https://api-na.myconnectwise.net/v4_6_release/apis/3.0"
)

type Client struct {
	httpClient   *http.Client
	encodedCreds string
	clientId     string
}

func NewClient(httpClient *http.Client, pubkey, privKey, clientId, companyId string) *Client {
	username := fmt.Sprintf("%s+%s", companyId, pubkey)
	return &Client{
		httpClient:   httpClient,
		encodedCreds: basicAuth(username, privKey),
		clientId:     clientId,
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *Client) request(ctx context.Context, method, endpoint string, payload io.Reader, target interface{}) error {
	u := fmt.Sprintf("%s/%s", baseUrl, endpoint)
	req, err := http.NewRequestWithContext(ctx, method, u, payload)
	if err != nil {
		return fmt.Errorf("creating the request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("clientId", c.clientId)
	req.Header.Set("Authorization", c.encodedCreds)

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
