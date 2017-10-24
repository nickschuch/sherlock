package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func notifySlack(key, channel, bucket, cluster, namespace, pod, container, incident string) error {
	api := slack.New(key)

	params := slack.PostMessageParameters{
		Username:  fmt.Sprintf("Watson - %s", cluster),
		IconEmoji: ":watson:",
	}

	_, _, err := api.PostMessage(channel, slackMessage(namespace, pod, container, bucket, incident), params)
	return err
}

// Helper function for generating a slack message.
func slackMessage(namespace, pod, container, bucket, incident string) string {
	return fmt.Sprintf("A Pod has been murdered: *%s* / *%s* / *%s*\n`sherlock inspect --s3-bucket=%s %s`", namespace, pod, container, bucket, incident)
}
