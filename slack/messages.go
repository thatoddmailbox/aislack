package slack

import "net/http"

type SentMessage struct {
	Text string `json:"text"`
	User string `json:"user"`
	TS   string `json:"ts"`
}

type historyResponse struct {
	OK       bool          `json:"ok"`
	Messages []SentMessage `json:"messages"`
}

func (c *Client) ListConversationHistory(channelID ChannelID) ([]SentMessage, error) {
	var response historyResponse
	err := c.request(http.MethodGet, "conversations.history", map[string]string{
		"channel": string(channelID),
	}, &response)
	if err != nil {
		return nil, err
	}

	return response.Messages, nil
}

func (c *Client) PostMessage(channelID ChannelID, text string) error {
	err := c.request(http.MethodPost, "chat.postMessage", map[string]string{
		"channel": string(channelID),
		"text":    text,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
