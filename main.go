package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/nickschuch/sherlock/cmd"
)

func main() {
	app := kingpin.New("Sherlock", "When a Pod is murdered, Sherlock isn't far away to solve the mystery.")

	cmd.Watson(app)
	cmd.Inspect(app)
	cmd.List(app)
	cmd.Dummy(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
