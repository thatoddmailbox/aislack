package slack

import "net/http"

type MessageTS string

type SentMessage struct {
	Text string    `json:"text"`
	User string    `json:"user"`
	TS   MessageTS `json:"ts"`
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

func (c *Client) ReactToMessage(emoji string, messageTS MessageTS, channelID ChannelID) error {
	err := c.request(http.MethodPost, "reactions.add", map[string]string{
		"channel":   string(channelID),
		"name":      emoji,
		"timestamp": string(messageTS),
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UnreactToMessage(emoji string, messageTS MessageTS, channelID ChannelID) error {
	err := c.request(http.MethodPost, "reactions.remove", map[string]string{
		"channel":   string(channelID),
		"name":      emoji,
		"timestamp": string(messageTS),
	}, nil)
	if err != nil {
		return err
	}

	return nil
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
