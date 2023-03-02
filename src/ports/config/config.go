package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type rootConfig struct {
	path   string
	config *models.Config
}

func NewConfigService() modelhelper.ConfigService {
	return &rootConfig{}
}
func Load() *models.Config {
	loader := NewConfigLoader()
	cfg, err := loader.Load()

	if err != nil {
		// handle error
	}

	return cfg
}

func NewConfigLoader() modelhelper.ConfigService {
	return &rootConfig{}
}

func New() *models.Config {

	usr, err := user.Current()
	if err != nil {

	}

	c := models.Config{
		Port:          3003,
		ConfigVersion: "3.0",
	}

	if usr != nil {
		c.Developer.Name = usr.Name
	}
	return &c
}

// Load returns a new default configuration
func (c *rootConfig) Load() (*models.Config, error) {
	path := filepath.Join(Location(), "config.yaml")

	cfg, err := c.LoadFromFile(path)
	if err != nil {
		return nil, err
	}

	cfg.DirectoryName = Location()
	return cfg, nil

}

// func (c *rootConfig) GetConnections() (*map[string]models.Connection, error) {
// 	return &c.config.Connections, nil
// }

func (c *rootConfig) LoadFromFile(path string) (*models.Config, error) {
	var cfg *models.Config

	dat, e := ioutil.ReadFile(path)
	if e != nil {
		log.Fatalf("cannot load file: %v", e)
		return nil, e
	}

	err := yaml.Unmarshal(dat, &cfg)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
		return nil, e
	}

	return cfg, nil
}

// ConfigFolder returns the root path of ModelHelper
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

func (cfg *rootConfig) Save() error {
	if _, err := os.Stat(cfg.path); os.IsNotExist(err) {
		os.Mkdir(cfg.path, 0755)
	}

	return update(cfg.config)
}
