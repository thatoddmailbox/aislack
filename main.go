package main

import (
	"log"
	"strings"
	"time"

	"github.com/thatoddmailbox/aislack/config"
	"github.com/thatoddmailbox/aislack/jobmgr"
	"github.com/thatoddmailbox/aislack/slack"
	"github.com/thatoddmailbox/aislack/slack/webhookserver"
)

var jobmgrClient *jobmgr.Client
var slackClient *slack.Client

func genimage(command string, text string, c slack.CommandContext) (map[string]string, error) {
	go func() {
		time.Sleep(1 * time.Second)

		messages, err := slackClient.ListConversationHistory(c.ChannelID)
		if err != nil {
			panic(err)
		}

		var commandMessage *slack.SentMessage
		lookingForText := strings.TrimSpace(command + " " + text)
		for _, message := range messages {
			if message.Text == lookingForText && message.User == c.UserID {
				commandMessage = &message
				break
			}
		}

		if commandMessage == nil {
			// couldn't find the message
			// TODO: handle this better? is there a better way to handle this?
			return
		}

		err = slackClient.ReactToMessage("thonk", commandMessage.TS, c.ChannelID)
		if err != nil {
			panic(err)
		}

		jobID, err := jobmgrClient.StartJob("sd-txt2img", map[string]string{
			"prompt":        text,
			"slackUserID":   c.UserID,
			"slackUserName": c.UserName,
		})
		if err != nil {
			panic(err)
		}

		for {
			time.Sleep(2 * time.Second)

			job, _, artifacts, err := jobmgrClient.GetJob(jobID)
			if err != nil {
				panic(err)
			}

			if job.Status == jobmgr.JobStatusCompleted {
				// it's completed!
				err = slackClient.ReactToMessage("white_check_mark", commandMessage.TS, c.ChannelID)
				if err != nil {
					panic(err)
				}

				err = slackClient.UnreactToMessage("thonk", commandMessage.TS, c.ChannelID)
				if err != nil {
					panic(err)
				}

				// upload the artifact
				artifact := artifacts[0]

				artifactReader, err := jobmgrClient.DownloadArtifact(&artifact)
				if err != nil {
					panic(err)
				}
				defer artifactReader.Close()

				err = slackClient.UploadFile(text, "\""+text+"\" - <@"+c.UserID+">", c.ChannelID, artifactReader)
				if err != nil {
					panic(err)
				}

				break
			} else if job.Status == jobmgr.JobStatusFailed {
				// it's failed!
				err = slackClient.ReactToMessage("cry", commandMessage.TS, c.ChannelID)
				if err != nil {
					panic(err)
				}

				err = slackClient.UnreactToMessage("thonk", commandMessage.TS, c.ChannelID)
				if err != nil {
					panic(err)
				}

				break
			}
		}
	}()

	return map[string]string{
		"response_type": "in_channel",
	}, nil
}

func main() {
	log.Println("aislack")

	err := config.Load()
	if err != nil {
		panic(err)
	}

	jobmgrClient, err = jobmgr.NewClient()
	if err != nil {
		panic(err)
	}

	slackClient, err = slack.NewClient(config.Current.Slack.UserOAuthToken)
	if err != nil {
		panic(err)
	}

	slackClient.RegisterCommandHandler("/genimage", genimage)
	slackClient.RegisterCommandHandler("/imagetest", genimage)

	whserver, err := webhookserver.NewServer(config.Current.WebhookServer.Port, slackClient)
	if err != nil {
		panic(err)
	}

	whserver.StartListening()
}
