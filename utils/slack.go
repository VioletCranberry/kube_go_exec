package utils

import (
	"fmt"

	"github.com/slack-go/slack"
)

type SlackClient struct {
	client *slack.Client
}

func InitSlackApi(token string) *SlackClient {
	var client = new(SlackClient)
	client.client = slack.New(token)
	return client
}

func (slack_client *SlackClient) SendToChannel(
	channel_id string, command string, pod string, stdout string, stderr string) (string, slack.Attachment, error) {
	var message, color string

	if stdout != "" {
		message = stdout
		color = "#04bf00"
	}
	if stderr != "" {
		message = stderr
		color = "#f8002f"
	}
	attachment := slack.Attachment{
		Color:   color,
		Title:   pod,
		Pretext: fmt.Sprintf("command: `%s`", command),
		Text:    fmt.Sprintf("```%s```", message),
	}
	channelid, _, err := slack_client.client.PostMessage(
		channel_id,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		return "", attachment, err
	}
	return channelid, attachment, err
}
