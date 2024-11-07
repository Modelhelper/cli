package config

import (
	"fmt"
	"log"
	"log/slog"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"gopkg.in/yaml.v3"
)

type rootConfig struct {
	path   string
	config *models.Config
}

// Create implements modelhelper.ConfigService.
func (c *rootConfig) Create() (*models.Config, error) {
	cfg := *&models.Config{}

	var (
		// codeEditor string
		dbTplPath      string
		codeTplPath    string
		projectTplPath string
		langDefPath    string
	)

	pathExists := func(path string) error {
		if len(path) == 0 {
			return nil
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("Path does not exist")
		} else {
			return nil
		}
	}

	editorSelect := huh.NewSelect[string]().
		Title("It looks like this is the first time the ModelHelper is run on this computer").
		Description("Please select the type of editor you want to connect use to view things in").
		Options(
			huh.NewOption("I do not want to set a default one", "none"),
			huh.NewOption("Visual Studio Code (code)", "code"),
			huh.NewOption("Vim", "vim"),
			huh.NewOption("Nano", "nano"),
		).
		Value(&cfg.DefaultEditor)

	tplDatabaseInput := huh.NewInput().
		Prompt("Database templates: ").
		Validate(pathExists).
		Value(&dbTplPath)
	tplProjectInput := huh.NewInput().
		Prompt("Project templates: ").
		Validate(pathExists).
		Value(&projectTplPath)

	tplCodeInput := huh.NewInput().
		Title("Add default template paths (leave blank if you want to add later)").
		Prompt("Code templates: ").
		Validate(pathExists).
		Value(&codeTplPath)

	// langDefInput := huh.NewInput().
	// 	Title("Where can we find language definition files").
	// 	Prompt("Path to lang def: ").
	// 	Validate(pathExists).
	// 	Value(&langDefPath)

	huh.NewForm(
		huh.NewGroup(editorSelect),
		huh.NewGroup(tplCodeInput, tplDatabaseInput, tplProjectInput),
		// huh.NewGroup(langDefInput),
	).Run()

	if len(codeTplPath) > 0 {
		cfg.Templates.Code = append(cfg.Templates.Code, codeTplPath)
	}

	if len(dbTplPath) > 0 {
		cfg.Templates.Database = append(cfg.Templates.Database, dbTplPath)
	}

	if len(projectTplPath) > 0 {
		cfg.Templates.Project = append(cfg.Templates.Project, projectTplPath)
	}

	if len(langDefPath) > 0 {
		cfg.Languages.Definitions = langDefPath
	}

	return &cfg, nil
}

// ConfigExists implements modelhelper.ConfigService.
func (c *rootConfig) ConfigExists() bool {
	path := filepath.Join(Location(), "config.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
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

	dat, e := os.ReadFile(path)
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
func (cfg *rootConfig) SaveConfig(c *models.Config) error {
	slog.Info("Saving config", "path", Location())
	if _, err := os.Stat(Location()); os.IsNotExist(err) {
		os.Mkdir(Location(), 0755)
	}

	return update(c)
}
