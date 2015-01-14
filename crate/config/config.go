// Handles the configuration for the crate package

package config

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//=============================================================================

type Config struct {
	Debug  bool     `yaml:debug,omitempty`  // default false
	Notify []string `yaml:notify,omitempty` // default []
	Level  string   `yaml:level,omitempty`  // default INFO
}

//=============================================================================

// Create a New Config with default values, one of two ways to create a config
func New() *Config {

	// Create the config struct
	config := new(Config)

	// Set the default values
	config.Debug = false
	config.Notify = make([]string, 0, 0)
	config.Level = "INFO"

	return config
}

// Load a Config by loading a YAML file, second of two ways to create a config
func Load(path string) (*Config, error) {

	// Check to ensure that the path exists
	if exists, _ := PathExists(path); !exists {
		return nil, errors.New("No YAML config file at specified path")
	}

	// Load the file from the path
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshall the YAML into a config object
	config := New()
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

//=============================================================================

// Dump a Config to a YAML file, primary way of saving a config to disk
func (conf *Config) Dump(path string) error {

	// Marshall the data to write to file
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)

}
