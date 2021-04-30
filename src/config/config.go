package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"modelhelper/cli/project"
	"modelhelper/cli/source"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {

	// ConfigVersion gets the version that this configuration file is using.
	ConfigVersion     string                       `json:"configVersion" yaml:"configVersion"`
	AppVersion        string                       `json:"appVersion" yaml:"appVersion"`
	Connections       map[string]source.Connection `json:"connections" yaml:"connections"`
	DefaultConnection string                       `json:"defaultConnection" yaml:"defaultConnection"`

	Templates struct {
		Location string `json:"configVersion" yaml:"configVersion"`
	} `json:"templates" yaml:"templates"`
	Languages struct {
		Definitions string `json:"definitions" yaml:"definitions"`
	} `json:"languages" yaml:"languages"`
	Logging struct {
		Enabled bool `json:"enabled" yaml:"enabled"`
	} `json:"logging" yaml:"logging"`
}

type LanguageDef struct {
	Version        string                     `json:"version" yaml:"version"`
	Language       string                     `json:"language" yaml:"language"`
	Datatypes      map[string]Datatype        `json:"datatypes" yaml:"datatypes"`
	Keys           map[string]project.CodeKey `json:"Keys" yaml:"Keys"`
	UsesNamespace  bool                       `json:"usesNamespace" yaml:"usesNamespace"`
	CanInject      bool                       `json:"canInject" yaml:"canInject"`
	DefaultImports []string                   `json:"defaultImports" yaml:"defaultImports"`
}

type Datatype struct {
	NotNull             string `json:"notNull" yaml:"notNull"`
	Nullable            string `json:"nullable" yaml:"nullable"`
	NullableAlternative string `json:"nullableAlternative" yaml:"nullableAlternative"`
}

// Load returns a new default configuration
func Load() *Config {
	path := filepath.Join(Location(), "config.yaml")
	return LoadFromFile(path)

}

// Initialize builds the configuration
func (c *Config) Initialize() error {

	fmt.Println("Initialize stuff from config here")

	return nil
}

func (c *Config) GetLanguageDefs() (*map[string]LanguageDef, error) {
	return nil, nil
}
func (c *Config) GetConnections() (*map[string]source.Connection, error) {
	return &c.Connections, nil
}

func LoadLanguageFile(path string) (*LanguageDef, error) {
	var lang *LanguageDef

	dat, e := ioutil.ReadFile(path)
	if e != nil {
		log.Fatalf("cannot load file: %v", e)
		return nil, e
	}

	err := yaml.Unmarshal(dat, &lang)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
		return nil, err
	}

	return lang, nil
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
