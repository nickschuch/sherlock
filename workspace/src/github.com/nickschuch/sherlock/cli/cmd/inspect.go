package cmd

import (
	"fmt"
	"strings"

	"github.com/nickschuch/str2color"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/nickschuch/sherlock/common/storage"
)

type cmdInspect struct {
	region   string
	bucket   string
	incident string
}

func (cmd *cmdInspect) run(c *kingpin.ParseContext) error {
	// Load the incidents.
	client := storage.New(cmd.region, cmd.bucket)

	files, err := client.IncidentDetails(cmd.incident)
	if err != nil {
		return fmt.Errorf("failed to lookup incident: %s", err)
	}

	for file, content := range files {
		for _, line := range strings.Split(content, "\n") {
			fmt.Println(str2color.Wrap(file), "\t", line)
		}

		fmt.Println("-----------------------------------------------------------------------------------------------")
	}

	return nil
}

// Inspect declares the "inspect" sub command.
func Inspect(app *kingpin.Application) {
	c := new(cmdInspect)

	cmd := app.Command("inspect", "Inspect the incident").Action(c.run)
	cmd.Flag("s3-region", "Region of the S3 bucket to store data").Default("ap-southeast-2").OverrideDefaultFromEnvar("SHERLOCK_S3_REGION").StringVar(&c.region)
	cmd.Flag("s3-bucket", "Name of the S3 bucket to store data").Default("").OverrideDefaultFromEnvar("SHERLOCK_S3_BUCKET").StringVar(&c.bucket)
	cmd.Arg("incident", "ID of the incident").Required().StringVar(&c.incident)
}
