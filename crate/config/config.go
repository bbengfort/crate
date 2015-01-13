// Handles the configuration for the crate package

package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"
	// "gopkg.in/yaml.v2"
)

var cratePath string // Internal config path identifier

//=============================================================================

// Helper function to check if a path exists on the file system
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

//=============================================================================

// Figures out the configuration and data path and returns it or an error
// Configuration and data are in the HOME directory of the current user
func CrateDirectory() (string, error) {

	// Cache the crate directory path
	if cratePath == "" {

		// Configure crate directory path with home directory and operating system
		if homePath, err := homedir.Dir(); err == nil {

			if runtime.GOOS == "windows" {
				// Currently untested since I don't have windows
				cratePath = filepath.Join(homePath, "AppData", "Roaming", "Crate")
			} else {
				cratePath = filepath.Join(homePath, ".crate")
			}

		} else {
			return "", err
		}

	}

	// Ensure that the crate directory path exists and is initialized
	if err := InitializeCrateDirectory(cratePath); err != nil {
		cratePath = ""
		return "", err
	}

	return cratePath, nil

}

// Creates the Crate directory and initializes it with default files
func InitializeCrateDirectory(path string) error {

	// Check if the path does not exist already
	if exists, err := PathExists(path); !exists {
		if err != nil {
			return err
		}

		// Create the crate directory at the specifed path
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}

	}

	return nil
}

//=============================================================================
