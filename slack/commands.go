package slack

import (
	"net/url"
)

type CommandContext struct {
	ResponseURL string
	ChannelID   ChannelID
	TeamID      string
	TriggerID   string
	UserID      string
	UserName    string
}

type CommandHandler func(command string, text string, c CommandContext) (map[string]string, error)

func (c *Client) HandleSlashCommand(data url.Values) (map[string]string, error) {
	c.commandLock.Lock()
	defer c.commandLock.Unlock()

	command := data["command"][0]
	text := data["text"][0]
	responseURL := data["response_url"][0]
	channelID := ChannelID(data["channel_id"][0])
	teamID := data["team_id"][0]
	triggerID := data["trigger_id"][0]
	userID := data["user_id"][0]
	userName := data["user_name"][0]

	handler, ok := c.commandMap[command]
	if !ok {
		return map[string]string{
			"text":          "Slash command not recognized. Probably whoever set up this bot did a bad job",
			"response_type": "in_channel",
		}, nil
	}

	result, err := handler(command, text, CommandContext{
		ResponseURL: responseURL,
		ChannelID:   channelID,
		TeamID:      teamID,
		TriggerID:   triggerID,
		UserID:      userID,
		UserName:    userName,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) RegisterCommandHandler(command string, handler CommandHandler) {
	c.commandLock.Lock()
	defer c.commandLock.Unlock()

	c.commandMap[command] = handler
}
