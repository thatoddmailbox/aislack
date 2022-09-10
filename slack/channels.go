package slack

import "net/http"

type ChannelID string

func (c *Client) ListChannels() error {
	err := c.request(http.MethodGet, "conversations.list", nil)
	if err != nil {
		return err
	}

	return nil
}
