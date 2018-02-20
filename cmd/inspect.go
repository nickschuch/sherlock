package cmd

import (
	"fmt"
	"strings"

	"github.com/nickschuch/str2color"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/pkg/errors"

	"github.com/nickschuch/sherlock/storage"
	"github.com/nickschuch/sherlock/storage/types"
	"github.com/nickschuch/sherlock/utils/highlight"
)

var keywords = []string{
	"OOMKiller",
	"ERROR",
	"Error",
	"error",
	"FATAL",
	"Fatal",
	"fatal",
	"unhealthy",
}

type cmdInspect struct {
	storage string
	region   string
	bucket   string
	incident string
}

func (cmd *cmdInspect) run(c *kingpin.ParseContext) error {
	client, err := storage.New(cmd.storage, cmd.region, cmd.bucket)
	if err != nil {
		return errors.Wrap(err, "failed to get client")
	}

	inspect, err := client.Inspect(types.InspectParams{
		ID: cmd.incident,
	})
	if err != nil {
		return errors.Wrap(err, "failed to lookup incident")
	}

	for _, clue := range inspect.Incident.Clues {
		for _, line := range strings.Split(clue.Content, "\n") {
			fmt.Println(str2color.Wrap(clue.Name), "\t", highlight.Wrap(line, keywords))
		}
	}

	return nil
}

// Inspect declares the "inspect" sub command.
func Inspect(app *kingpin.Application) {
	c := new(cmdInspect)

	cmd := app.Command("inspect", "Inspect the incident").Action(c.run)
	cmd.Flag("storage", "Type of storage which has our incidents").Default("s3").Envar("SHERLOCK_STORAGE").StringVar(&c.storage)
	cmd.Flag("region", "Region of the S3 bucket to store data").Default("ap-southeast-2").Envar("SHERLOCK_REGION").StringVar(&c.region)
	cmd.Flag("bucket", "Name of the S3 bucket to store data").Required().Envar("SHERLOCK_BUCKET").StringVar(&c.bucket)
	cmd.Arg("id", "ID of the incident").Required().StringVar(&c.incident)
}
