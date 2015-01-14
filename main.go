// Command-Line program and main package for the Crate Daemon

package main

import (
	"os"

	"github.com/bbengfort/crate/crate"
	"github.com/bbengfort/crate/crate/version"
	"github.com/codegangsta/cli"
)

var console *crate.Console

func main() {

	app := cli.NewApp()
	app.Name = "crate"
	app.Usage = "file archival and metadata synchronization tool (experimental)"
	app.Version = version.Version()
	app.Author = "Benjamin Bengfort"
	app.Email = "benjamin@bengfort.com"
	app.Flags = []cli.Flag{
		cli.BoolFlag{"debug", "set debug mode for vebose logging", ""},
	}

	app.Action = func(c *cli.Context) {

		// debug := c.Bool("debug")

		service := new(crate.CrateService)
		service.Backup(c.Args()[0])

	}

	app.Run(os.Args)

}
