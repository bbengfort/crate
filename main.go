// Command-Line program and main package for the Crate Daemon

package main

import (
	"github.com/bbengfort/crate/crate"
	"github.com/bbengfort/crate/version"
	"path/filepath"
)

func main() {

	// Create debug console for now
	console := new(crate.Console)
	console.Init(true)
	console.Info("Crate version %s", version.Version())

	path, _ := crate.NewPath("../benfs")

	if fm, ok := path.(*crate.FileMeta); ok {
		console.Info("Path is %s, isdir: %v, isfile: %v", fm, fm.IsDir(), fm.IsFile())
		console.Info("Base: %s Dir: %s, Ext: %s", fm.Base(), fm.Dir(), fm.Ext())
	} else if dir, ok := path.(*crate.Dir); ok {
		console.Info("Path is %s, isdir: %v, isfile: %v", dir, dir.IsDir(), dir.IsFile())
		console.Info("Dir: %s", dir.Dir())
		console.Info("Subdirectory = %s", dir.Join("path", "to", "a", "subdir"))

		dir.Walk(func(path crate.Path, err error) error {

			if err != nil {
				return err
			}

			// if !path.IsHidden() {
			// 	console.Info("%s", path)
			// }

			if !path.Dir().IsHidden() {
				console.Info("%s", path.Dir())
				return filepath.SkipDir
			}

			return nil

		})

	}
}
