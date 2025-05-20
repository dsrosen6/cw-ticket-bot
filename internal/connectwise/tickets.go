package connectwise

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) GetTicketNote(ctx context.Context, ticketId int) ([]Note, error) {
	q := url.QueryEscape("_info/dateEntered desc")
	endpoint := fmt.Sprintf("service/tickets/%d/allNotes?orderBy=%s", ticketId, q)
	var n []Note

	if err := c.request(ctx, "GET", endpoint, nil, &n); err != nil {
		return nil, fmt.Errorf("getting the ticket notes: %w", err)
	}

	return n, nil
}
