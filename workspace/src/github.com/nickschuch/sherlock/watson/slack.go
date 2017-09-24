package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func notifySlack(key, channel, bucket, cluster, namespace, pod, container, incident string) error {
	api := slack.New(key)

	params := slack.PostMessageParameters{
		Username:  "Watson",
		IconEmoji: ":watson:",
		Attachments: []slack.Attachment{
			{
				Fields: []slack.AttachmentField{
					{
						Title: "Cluster",
						Value: cluster,
						Short: true,
					},
					{
						Title: "Pod",
						Value: fmt.Sprintf("%s / %s / %s", namespace, pod, container),
					},
					{
						Title: "Inspect",
						Value: fmt.Sprintf("sherlock inspect --s3-bucket=%s %s", bucket, incident),
					},
				},
			},
		},
	}

	_, _, err := api.PostMessage(channel, "A Pod has been murdered!", params)
	return err
}
