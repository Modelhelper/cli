package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

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
