package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/nickschuch/sherlock/cli/cmd"
)

func main() {
	app := kingpin.New("Sherlock", "When a Pod is murdered, Sherlock isn't far away to solve the mystery.")

	// Setup all the subcommands.
	cmd.List(app)
	cmd.Inspect(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
