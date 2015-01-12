// Command-Line program and main package for the Crate Daemon

package main

import (
	"os"
	"path/filepath"

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

		debug := c.Bool("debug")

		// Create debug console for now
		console = new(crate.Console)
		console.Init(debug)
		crate.InitMagic()
		defer crate.Magic.Close()

		path, err := crate.NewPath(c.Args()[0])
		if err != nil {
			console.Fatal("Could not open path: %s (%s)", c.Args()[0], err)
		}

		if fm, ok := path.(*crate.FileMeta); ok {

			if img, ok := crate.ConvertImageMeta(fm); ok {
				console.Log(img.Info())
			} else {
				console.Log(fm.Info())
			}

		} else if dir, ok := path.(*crate.Dir); ok {
			// Otherwise walk the directory for stats about it

			files := make(map[string]int)
			tags := make(map[string]int)

			dir.Walk(func(path crate.Path, err error) error {
				if err != nil {
					return err
				}

				if path.Dir().IsHidden() {
					console.Info("testing dir %s hidden on %s", path.Dir(), path)
					return filepath.SkipDir
				}

				if path.IsFile() && !path.IsHidden() {
					mtype, _ := crate.MimeType(path.String())
					files[mtype] += 1

					if img, ok := crate.ConvertImageMeta(path.(*crate.FileMeta)); ok {
						if img.IsJPEG() {
							img.Populate()
							for key, _ := range img.Tags {
								tags[key] += 1
							}
						}
					}
				}

				return nil

			})

			console.Log("File Mimetypes")
			for k, v := range files {
				console.Log("%d: %s", v, k)
			}

			console.Log("\n\nJPEG Tags")
			for k, v := range tags {
				console.Log("%d: %s", v, k)
			}

		}
	}

	app.Run(os.Args)

}
