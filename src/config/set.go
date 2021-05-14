package config

import (
	"io/ioutil"
	"modelhelper/cli/source"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func SetConnection(key string, c *source.Connection, isDefault bool, merge bool) error {
	cfg := Load()

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
func SetDeveloper(name string, email string) error {
	cfg := Load()
	dev := Developer{name, email}

	cfg.Developer = dev
	return update(cfg)
}
func SetPort(api int, web int) error {
	cfg := Load()

	cfg.ApiPort = api
	cfg.WebPort = web

	return update(cfg)
}

func SetTemplateLocation(loc string) error {
	cfg := Load()
	cfg.Templates.Location = loc

	return update(cfg)
}
func SetLangDefLocation(loc string) error {
	cfg := Load()
	cfg.Languages.Definitions = loc

	return update(cfg)
}

func update(cfg *Config) error {

	d, err := yaml.Marshal(&cfg)

	if err != nil {

		return err
	}

	path := filepath.Join(Location(), "config.yaml")
	err = ioutil.WriteFile(path, d, 0777)

	return err

}
