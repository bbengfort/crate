// Command-Line program and main package for the Crate Daemon

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bbengfort/crate/crate"
	"github.com/bbengfort/crate/crate/version"
)

var console *crate.Console

// Print the Usage of crate to Stderr
var Usage = func() {
	fmt.Fprintln(os.Stderr, "usage: crate [-h] [-v] directory\n")
	fmt.Fprintln(os.Stderr, "crate watches a directory for changes\n")
	fmt.Fprintln(os.Stderr, "optional arguments:")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage

	// Add command line arguments
	help := flag.Bool("h", false, "print help and exit")
	vers := flag.Bool("v", false, "print version and exit")
	debug := flag.Bool("debug", false, "set debug mode for verbose logging")

	// Parse the command line arguments
	flag.Parse()

	// Create debug console for now
	console = new(crate.Console)
	console.Init(*debug)
	crate.InitMagic()
	defer crate.Magic.Close()

	// Handle help, version, and arguments
	if *help {
		Usage()
		os.Exit(0)
	} else if *vers {
		console.Log("Crate version %s", version.Version())
		os.Exit(0)
	} else if flag.NArg() != 1 {
		Usage()
		os.Exit(2)
	}

	path, err := crate.NewPath(flag.Arg(0))
	if err != nil {
		console.Fatal("Could not open path: %s (%s)", flag.Arg(0), err)
	}

	if fm, ok := path.(*crate.FileMeta); ok {
		// Handle a file passed (just report the mimetype)
		mtype, err := crate.MimeType(fm.Path)
		if err != nil {
			console.Fatal("Could not magic type: %s", err)
		}

		console.Log("Path %s is a %s", fm, mtype)

	} else if dir, ok := path.(*crate.Dir); ok {
		// Otherwise walk the directory for stats about it

		files := make(map[string]int)

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
			}

			return nil

		})

		for k, v := range files {
			console.Log("%d: %s", v, k)
		}

	}
}
