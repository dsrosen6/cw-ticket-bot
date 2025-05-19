package webex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type Message struct {
	RoomId string `json:"roomId,omitempty"`
	Person string `json:"toPersonEmail,omitempty"`
	Text   string `json:"markdown"`
}

func NewMessageToPerson(email, text string) Message {
	return Message{Person: email, Text: text}
}

func NewMessageToRoom(roomId, text string) Message {
	return Message{RoomId: roomId, Text: text}
}

func (c *Client) SendMessage(ctx context.Context, message Message) error {
	j, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshaling message to json: %w", err)
	}

	p := bytes.NewReader(j)

	if err := c.request(ctx, "POST", "messages", p, nil); err != nil {
		return err
	}

	return nil
}
