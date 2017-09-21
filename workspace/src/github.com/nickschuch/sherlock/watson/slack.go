package main

import (
	"bytes"
	"text/template"

	slack "github.com/nickschuch/go-slack"
)

const messageFormat = `A Pod murder has taken place!!!
**Cluster**: _{{ .Cluster }}_
**Details**: _{{ .Name }}_
**Inspect**: _sherlock --bucket=foo inspect {{ .Incident }}_`

func notifySlack(slackUrl, slackEmoji, cluster, name, incident string) error {
	msg, err := notifySlackMessage(cluster, name, incident)
	if err != nil {
		return err
	}

	return slack.Send("Watson", slackEmoji, msg, slackUrl)
}

// Helper function for building hostname.
func notifySlackMessage(cluster, name, incident string) (string, error) {
	var formatted bytes.Buffer

	msg := Message{
		Cluster:  cluster,
		Name:     name,
		Incident: incident,
	}

	err := template.Must(template.New("message").Parse(messageFormat)).Execute(&formatted, msg)
	if err != nil {
		return "", err
	}

	return formatted.String(), nil
}
