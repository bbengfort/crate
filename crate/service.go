// Implements the primary service functionality of crate

package crate

import (
	"path/filepath"

	"github.com/bbengfort/crate/crate/config"
)

type CrateService struct {
	initialized bool           // Whether or not the service is initialized
	conf        *config.Config // Stores the configuration of the service
}

//=============================================================================

// Initialize all the various services, remember to also defer close!
func (service *CrateService) Init() {
	// Initialize the console
	InitializeConsole()

	// Initialize the configuration
	path, err := config.CrateConfigPath()
	if err != nil {
		console.Fatal("Could not initialize configuration: %s", err)
	}
	service.conf, err = config.Load(path)
	if err != nil {
		console.Fatal("Could not load configuration: %s", err)
	}

	// Initialize the logger
	if err := InitializeLoggers(LevelFromString(service.conf.Level)); err != nil {
		console.Fatal("Could not initialize loggers: %s", err)
	}

	// Initialize the database
	if err := InitializeDatabase(); err != nil {
		console.Fatal("Could not initialize the database: %s", err)
	}

	// Initialize the libmagic database
	if err := InitMagic(); err != nil {
		console.Fatal("Could not initialize libmagic: %s", err)
	}

	service.initialized = true
}

// Runs the backup utility service on a specified directory
func (service *CrateService) Backup(dirPath string) {
	if !service.initialized {
		service.Init()
	}

	// Defer closing of various utilities
	defer CloseDatabase()
	defer CloseLoggers()
	defer Magic.Close()

	rootPath, err := NewPath(dirPath)
	if err != nil {
		console.Fatal("Could not open path \"%s\": %s", dirPath, err)
	}

	if root, ok := rootPath.(*Dir); ok {

		// Primary functionality of Backup analysis
		// First log starting of backup on directory
		eventLogger.Info("started backup on directory \"%s\"", root)

		root.Walk(func(path Path, err error) error {

			// If there is an error, stop the walk
			if err != nil {
				return err
			}

			// Skip any files inside of hidden directories
			if path.Dir().IsHidden() {
				return filepath.SkipDir
			}

			// If this is a file (not a dir) and is not hidden
			if path.IsFile() && !path.IsHidden() {

				if fm, ok := path.(*FileMeta); ok {

					if img, ok := ConvertImageMeta(fm); ok {
						img.Store()
					} else {
						fm.Store()
					}

				} else {
					eventLogger.Error("could not convert \"%s\" to file meta object", path)
				}

			}

			return nil

		})

		// Log the completion of the backup on directory
		eventLogger.Info("finished backup on directory \"%s\"", root)

	} else {
		console.Fatal("Specified path is not a directory, \"%s\"", dirPath)
	}

}
