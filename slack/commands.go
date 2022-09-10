package slack

import "net/url"

type CommandHandler func(text string) (map[string]string, error)

func (c *Client) HandleSlashCommand(data url.Values) (map[string]string, error) {
	c.commandLock.Lock()
	defer c.commandLock.Unlock()

	command := data["command"][0]
	text := data["text"][0]

	handler, ok := c.commandMap[command]
	if !ok {
		return map[string]string{
			"text":          "Slash command not recognized. Probably whoever set up this bot did a bad job",
			"response_type": "in_channel",
		}, nil
	}

	result, err := handler(text)
	if err != nil {
		return nil, err
	}

	return result, nil
}
