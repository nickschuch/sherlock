package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/nickschuch/sherlock/common/storage"
)

type cmdList struct {
	region string
	bucket string
}

func (cmd *cmdList) run(c *kingpin.ParseContext) error {
	client := storage.New(cmd.region, cmd.bucket)

	incidents, err := client.Incidents()
	if err != nil {
		return err
	}

	table := uitable.New()
	table.MaxColWidth = 50

	table.AddRow("INCIDENT ID", "TIMESTAMP", "NAMESPACE", "POD", "CONTAINER")
	for id, incident := range incidents {
		table.AddRow(id, incident.Created, incident.Namespace, incident.Pod, incident.Container)
	}

	fmt.Println(table)

	return nil
}

// List declares the "list" sub command.
func List(app *kingpin.Application) {
	c := new(cmdList)

	cmd := app.Command("list", "List all incidents").Action(c.run)
	cmd.Flag("s3-region", "Region of the S3 bucket to store data").Default("ap-southeast-2").OverrideDefaultFromEnvar("SHERLOCK_S3_REGION").StringVar(&c.region)
	cmd.Flag("s3-bucket", "Name of the S3 bucket to store data").Default("").OverrideDefaultFromEnvar("SHERLOCK_S3_BUCKET").StringVar(&c.bucket)
}
