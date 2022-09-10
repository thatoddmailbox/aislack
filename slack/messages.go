package slack

import "net/http"

func (c *Client) PostMessage(channelID ChannelID, text string) error {
	err := c.request(http.MethodPost, "chat.postMessage", map[string]string{
		"channel": string(channelID),
		"text":    text,
	})
	if err != nil {
		return err
	}

	return nil
}
