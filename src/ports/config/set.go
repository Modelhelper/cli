package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"modelhelper/cli/modelhelper/models"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func SetDefaultConnection(key string) error {
	loader := NewConfigLoader()
	cfg, err := loader.Load()

	if err != nil {
		// handle error
	}
	_, f := cfg.Connections[key]

	if !f {
		m := fmt.Sprintf("The connection: %s does not exists and cannot be a default connection", key)
		return errors.New(m)
	}

	cfg.DefaultConnection = key

	return update(cfg)
}

func SetDefaultEditor(editor string) error {
	loader := NewConfigLoader()
	cfg, err := loader.Load()

	if err != nil {
		// handle error
	}

	cfg.DefaultEditor = strings.ToLower(editor)

	return update(cfg)
}

func SetConnection(key string, c *models.Connection, isDefault bool, merge bool) error {
	loader := NewConfigLoader()
	cfg, err := loader.Load()

	if err != nil {
		// handle error
	}

	current, exists := cfg.Connections[key]

	if exists && merge {
		if len(c.ConnectionString) > 0 {
			current.ConnectionString = c.ConnectionString
		}

		if len(c.Description) > 0 {
			current.Description = c.Description
		}
		if len(c.Type) > 0 {
			current.Type = c.Type
		}
		if len(c.Schema) > 0 {
			current.Schema = c.Schema
		}

	} else {
		current = *c
	}

	cfg.Connections[key] = current

	if isDefault {
		cfg.DefaultConnection = key
	}
	return update(cfg)
}
func SetDeveloper(name string, email string, github string, merge bool) error {
	loader := NewConfigLoader()
	cfg, err := loader.Load()

	if err != nil {
		// handle error
	}

	if merge {
		if len(name) > 0 {
			cfg.Developer.Name = name
		}

		if len(email) > 0 {
			cfg.Developer.Email = email
		}

		if len(github) > 0 {
			cfg.Developer.GitHubAccount = github
		}
	} else {

		dev := models.Developer{name, email, github}
		cfg.Developer = dev
	}

	return update(cfg)
}
func SetPort(api int, web int) error {
	loader := NewConfigLoader()
	cfg, err := loader.Load()

	if err != nil {
		// handle error
	}

	cfg.Port = api

	return update(cfg)
}

func SetTemplateLocation(loc string) error {
	loader := NewConfigLoader()
	cfg, err := loader.Load()

	if err != nil {
		// handle error
	}
	cfg.Templates.Code = append(cfg.Templates.Code, loc)

	return update(cfg)
}
func SetLangDefLocation(loc string) error {
	loader := NewConfigLoader()
	cfg, err := loader.Load()

	if err != nil {
		// handle error
	}
	cfg.Languages.Definitions = loc

	return update(cfg)
}

func update(cfg *models.Config) error {

	d, err := yaml.Marshal(&cfg)

	if err != nil {

		return err
	}

	path := filepath.Join(Location(), "config.yaml")
	err = ioutil.WriteFile(path, d, 0777)

	return err

}
