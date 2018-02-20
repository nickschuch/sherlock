package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/pkg/errors"

	"github.com/nickschuch/sherlock/storage"
	"github.com/nickschuch/sherlock/storage/types"
)

type cmdList struct {
	storage string
	region   string
	bucket   string
}

func (cmd *cmdList) run(c *kingpin.ParseContext) error {
	client, err := storage.New(cmd.storage, cmd.region, cmd.bucket)
	if err != nil {
		return errors.Wrap(err, "failed to get client")
	}

	list, err := client.List(types.ListParams{})
	if err != nil {
		return err
	}

	table := uitable.New()
	table.MaxColWidth = 50

	table.AddRow("ID", "CLUSTER", "TIMESTAMP", "NAMESPACE", "POD", "CONTAINER")
	for _, incident := range list.Incidents {
		table.AddRow(incident.ID, incident.Cluster, incident.Created, incident.Namespace, incident.Pod, incident.Container)
	}

	fmt.Println(table)

	return nil
}

// List declares the "list" sub command.
func List(app *kingpin.Application) {
	c := new(cmdList)

	cmd := app.Command("list", "List all incidents").Action(c.run)
	cmd.Flag("storage", "Type of storage which has our incidents").Default("s3").Envar("SHERLOCK_STORAGE").StringVar(&c.storage)
	cmd.Flag("region", "Region of the S3 bucket to store data").Default("ap-southeast-2").Envar("SHERLOCK_REGION").StringVar(&c.region)
	cmd.Flag("bucket", "Name of the S3 bucket to store data").Required().Envar("SHERLOCK_BUCKET").StringVar(&c.bucket)
}
