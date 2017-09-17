package cmd

import (
	"fmt"
	"os"
	"text/template"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/nickschuch/sherlock/common/storage"
)

const tmplInspect = `{{ range $key, $value := . }}
###########################################################################
{{ $key }}
###########################################################################

{{ $value }}
{{ end }}`

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

	return template.Must(template.New("inspect").Parse(tmplInspect)).Execute(os.Stdout, files)
}

// Inspect declares the "inspect" sub command.
func Inspect(app *kingpin.Application) {
	c := new(cmdInspect)

	cmd := app.Command("inspect", "Inspect the incident").Action(c.run)
	cmd.Flag("s3-region", "Region of the S3 bucket to store data").Default("ap-southeast-2").OverrideDefaultFromEnvar("SHERLOCK_S3_REGION").StringVar(&c.region)
	cmd.Flag("s3-bucket", "Name of the S3 bucket to store data").Default("").OverrideDefaultFromEnvar("SHERLOCK_S3_BUCKET").StringVar(&c.bucket)
	cmd.Arg("incident", "ID of the incident").Required().StringVar(&c.incident)
}
