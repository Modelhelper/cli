package config

import (
	"fmt"
	"modelhelper/cli/source"
	"os"
)

type Config struct {

	// ConfigVersion gets the version that this configuration file is using.
	ConfigVersion string                   `json:"configVersion"`
	AppVersion    string                   `json:"appVersion"`
	Sources       map[string]source.Source `json:"sources"`
	DefaultSource string                   `json:"defaultSource"`

	Templates struct {
		Location string `json:"configVersion"`
	} `json:"templates"`
	Languages struct {
		Definitions string `json:"definitions"`
	} `json:"languages"`
	Logging struct {
		Enabled bool `json:"enabled"`
	} `json:"logging"`
}

type LanguageDef struct {
	Definitions string
}

// New returns a new default configuration
func New() *Config {
	return nil
}

// Initialize builds the configuration
func (c *Config) Initialize() error {
	return nil
}

//ConfigFolder returns the root path of ModelHelper
func Location() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return fmt.Sprintf("%s/.modelhelper", homeDir)

}

// LocationExists checks if the config folder exists
func LocationExists() bool {
	homeDir := Location()

	if _, err := os.Stat(homeDir); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
