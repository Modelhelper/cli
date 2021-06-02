package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"modelhelper/cli/project"
	"modelhelper/cli/source"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {

	// ConfigVersion gets the version that this configuration file is using.
	ConfigVersion     string                       `json:"configVersion" yaml:"configVersion"`
	AppVersion        string                       `json:"appVersion" yaml:"appVersion"`
	Connections       map[string]source.Connection `json:"connections" yaml:"connections"`
	DefaultConnection string                       `json:"defaultConnection" yaml:"defaultConnection"`
	DefaultEditor     string                       `json:"editor" yaml:"editor"`
	Developer         Developer                    `json:"developer" yaml:"developer"`
	Port              int                          `json:"port" yaml:"port"`
	Code              project.ProjectCode          `json:"code" yaml:"code"`
	Templates         struct {
		Location string `json:"location" yaml:"location"`
	} `json:"templates" yaml:"templates"`
	Languages struct {
		Definitions string `json:"definitions" yaml:"definitions"`
	} `json:"languages" yaml:"languages"`
	Logging struct {
		Enabled bool `json:"enabled" yaml:"enabled"`
	} `json:"logging" yaml:"logging"`
}

type Developer struct {
	Name          string `json:"name" yaml:"name"`
	Email         string `json:"email" yaml:"email"`
	GitHubAccount string `json:"github" yaml:"github"`
}

func New() *Config {

	usr, err := user.Current()
	if err != nil {

	}

	c := Config{
		Port:          3003,
		ConfigVersion: "3.0",
	}

	if usr != nil {
		c.Developer.Name = usr.Name
	}
	return &c
}

// Load returns a new default configuration
func Load() *Config {
	path := filepath.Join(Location(), "config.yaml")
	return LoadFromFile(path)

}

func (c *Config) GetConnections() (*map[string]source.Connection, error) {
	return &c.Connections, nil
}

func LoadFromFile(path string) *Config {
	var cfg *Config

	dat, e := ioutil.ReadFile(path)
	if e != nil {
		log.Fatalf("cannot load file: %v", e)
		return nil
	}

	err := yaml.Unmarshal(dat, &cfg)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
		return nil
	}

	return cfg
}

//ConfigFolder returns the root path of ModelHelper
func Location() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return filepath.Join(homeDir, ".modelhelper")

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

func (cfg *Config) Save(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	return update(cfg)
}
