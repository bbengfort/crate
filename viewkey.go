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
	app.Name = "viewkey"
	app.Usage = "view a key in the database"
	app.Version = version.Version()
	app.Author = "Benjamin Bengfort"
	app.Email = "benjamin@bengfort.com"
	app.Flags = []cli.Flag{
		cli.BoolFlag{"debug", "set debug mode for vebose logging", ""},
		cli.BoolFlag{"all", "show all the database keys", ""},
		cli.IntFlag{"limit", 100, "limit the keys returned in all", ""},
	}

	app.Action = func(c *cli.Context) {

		debug := c.Bool("debug")
		all := c.Bool("all")
		limit := c.Int("limit")

		// Create debug console for now
		console = new(crate.Console)
		console.Init(debug)
		crate.InitMagic()
		defer crate.Magic.Close()

		// Connect to the database
		crate.InitializeDatabase()
		defer crate.CloseDatabase()

		if all {
			// TODO show all the keys in the database
			for _, key := range crate.FetchKeys(limit) {
				console.Log(key)
			}

		} else {
			for _, key := range c.Args() {
				path, err := crate.Fetch(key)
				if err != nil {
					console.Err("database lookup err", err)
				} else {
					console.Log(path.Info() + "\n")
				}
			}
		}
	}

	app.Run(os.Args)

}
