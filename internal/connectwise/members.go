package connectwise

import (
	"context"
	"fmt"
	"log/slog"
)

func (c *Client) GetMemberByIdentifier(ctx context.Context, identifier string) (*Member, error) {
	endpoint := fmt.Sprintf("system/members?conditions=identifier=\"%s\"", identifier)
	var m []Member

	if err := c.request(ctx, "GET", endpoint, nil, &m); err != nil {
		return nil, fmt.Errorf("getting member by identifier string: %w", err)
	}

	if len(m) == 0 {
		slog.Warn("member not found", "identifier", identifier)
		return nil, fmt.Errorf("member not found")
	}

	if len(m) > 1 {
		slog.Warn("multiple members found", "identifier", identifier)
		return nil, fmt.Errorf("multiple members found")
	}

	return &m[0], nil
}

func (c *Client) GetMember(ctx context.Context, memberId int) (*Member, error) {
	endpoint := fmt.Sprintf("system/members/%d", memberId)
	m := &Member{}

	if err := c.request(ctx, "GET", endpoint, nil, &m); err != nil {
		return nil, fmt.Errorf("getting member by id number: %w", err)
	}

	return m, nil
}
