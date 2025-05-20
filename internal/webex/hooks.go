package webex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) PostWebhook(ctx context.Context, newWebhook *Webhook) (*Webhook, error) {
	j, err := json.Marshal(newWebhook)
	if err != nil {
		return nil, fmt.Errorf("marshaling new webhook to json: %w", err)
	}

	p := bytes.NewReader(j)
	w := &Webhook{}
	if err := c.request(ctx, "POST", "webhooks", p, w); err != nil {
		return nil, fmt.Errorf("posting new webhook: %w", err)
	}

	return w, nil
}

func (c *Client) GetWebhooks(ctx context.Context) ([]Webhook, error) {
	w := &WebhooksGetResponse{}
	if err := c.request(ctx, "GET", "webhooks", nil, w); err != nil {
		return nil, fmt.Errorf("getting all webex webhooks: %w", err)
	}

	return w.Items, nil
}
