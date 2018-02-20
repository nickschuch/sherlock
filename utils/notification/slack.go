package notification

import (
	"fmt"

	"github.com/nlopes/slack"
)

type SlackParams struct {
	Key string
	Channel string
	Bucket string
	Cluster string
	Namespace string
	Pod string
	Container string
	ID string
}

func Slack(params SlackParams) error {
	api := slack.New(params.Key)

	post := slack.PostMessageParameters{
		Username:  fmt.Sprintf("Watson | %s", params.Cluster),
		IconEmoji: ":watson:",
		Attachments: []slack.Attachment{
			{
				Title: fmt.Sprintf("%s / %s / %s", params.Namespace, params.Pod, params.Container),
				Text:  fmt.Sprintf("sherlock inspect --bucket=%s %s", params.Bucket, params.ID),
			},
		},
	}

	_, _, err := api.PostMessage(params.Channel, "A pod has been murdered!", post)
	return err
}