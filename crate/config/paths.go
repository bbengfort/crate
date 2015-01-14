// Handles the configuration for the crate package

package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"
	// "gopkg.in/yaml.v2"
)

const (
	Windows          = "windows"
	WindowsAppData   = "AppData"
	WindowsRoaming   = "Roaming"
	WindowsCrateName = "Crate"
	UnixCrateName    = ".crate"
	DatabaseName     = "filemeta.db"
	ConfigName       = "config.yaml"
	LogDirName       = "logs"
	LogFileName      = "events.log"
)

var (
	cratePath   string // Internal config path identifier
	crateDBPath string // The path to the database storing the metadata
	configPath  string // The path to the YAML configuration file
	loggingPath string // The path to store the log files
)

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

// Helper function to clear the cache to force lookup
func ClearPathCache() {
	cratePath = ""
	crateDBPath = ""
	configPath = ""
	loggingPath = ""
}

//=============================================================================

// Figures out the configuration and data path and returns it or an error
// Configuration and data are in the HOME directory of the current user
func CrateDirectory() (string, error) {

	// Cache the crate directory path
	if cratePath == "" {

		// Configure crate directory path with home directory and operating system
		if homePath, err := homedir.Dir(); err == nil {

			if runtime.GOOS == Windows {
				// Currently untested since I don't have windows
				cratePath = filepath.Join(homePath, WindowsAppData, WindowsRoaming, WindowsCrateName)
			} else {
				cratePath = filepath.Join(homePath, UnixCrateName)
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

// Returns the database path and performs initialization if not done.
func CrateDatabasePath() (string, error) {

	// Ensure that there is a cratePath instantiated
	if cratePath == "" {
		if _, err := CrateDirectory(); err != nil {
			return "", err
		}
	}

	// Cache the the crate database path
	if crateDBPath == "" {
		crateDBPath = filepath.Join(cratePath, DatabaseName)
	}

	return crateDBPath, nil
}

// Returns the config path and performs initialization if not done
func CrateConfigPath() (string, error) {

	// Ensure that there is a cratePath instantiated
	if cratePath == "" {
		if _, err := CrateDirectory(); err != nil {
			return "", err
		}
	}

	// Cache the crate config path
	if configPath == "" {
		configPath = filepath.Join(cratePath, ConfigName)
		if err := InitializeConfigFile(configPath); err != nil {
			configPath = ""
			return "", err
		}
	}

	return configPath, nil
}

// Returns the logging path and performs initialization if not exists
func CrateLoggingPath() (string, error) {

	// Ensure that there is a cratePath instantiated
	if cratePath == "" {
		if _, err := CrateDirectory(); err != nil {
			return "", nil
		}
	}

	// Cache the crate logging path
	if loggingPath == "" {
		loggingDir := filepath.Join(cratePath, LogDirName)

		// Ensure that the crate directory path exists and is initialized
		if err := InitializeCrateDirectory(loggingDir); err != nil {
			return "", err
		}

		// Set the logging file in the logging dir
		loggingPath = filepath.Join(loggingDir, LogFileName)

	}

	return loggingPath, nil

}

//=============================================================================

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

// Creates a Config YAML file with the default settings at the specified path
func InitializeConfigFile(path string) error {

	// Check if the path exists already
	exists, err := PathExists(path)
	if err != nil {
		return err
	}

	// Don't overwrite existing config file
	if exists {
		return nil
	}

	// Create a new default config and write to path
	conf := New()
	if err := conf.Dump(path); err != nil {
		return err
	}

	return nil
}

//=============================================================================
