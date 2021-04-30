package project

import (
	"io/ioutil"
	"log"
	"modelhelper/cli/source"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Project struct {
	Version           string                       `yaml:"version"`
	Name              string                       `yaml:"name"`
	Language          string                       `yaml:"language"`
	DefaultConnection string                       `yaml:"defaultSource"`
	DefaultKey        string                       `yaml:"defaultKey"`
	Connections       map[string]source.Connection `yaml:"connections"`
	Code              ProjectCode                  `yaml:"code"`
	CustomerName      string                       `yaml:"customerName"`
	Header            string                       `yaml:"header"`
	Options           map[string]string            `yaml:"options"`
	Custom            interface{}                  `yaml:"custom"`
}

func DefaultLocation() string {
	p, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return filepath.Join(p, ".modelhelper", "project.yaml")
}

func Exists(path string) bool {

	pathInfo, err := os.Stat(path)

	if os.IsNotExist(err) || pathInfo.IsDir() {
		return false
	}

	return true
}

func Load(path string) (*Project, error) {
	if len(path) > 0 {
		pathInfo, err := os.Stat(path)
		if os.IsNotExist(err) || pathInfo.IsDir() {
			log.Fatal("Project does not exits")
			return nil, err
		}

	} else {
		p, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		path = filepath.Join(p, ".modelhelper", "project.yaml")
	}

	f, err := loadProjectFromFile(path)

	return f, err

}

func loadProjectFromFile(fileName string) (*Project, error) {
	var p *Project

	dat, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Fatalf("cannot load file: %v", e)
		return nil, e
	}

	err := yaml.Unmarshal(dat, &p)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
		return nil, err
	}

	return p, nil
}

func (p *Project) GetConnections() (*map[string]source.Connection, error) {
	return &p.Connections, nil
}

type ProjectCode struct {
	OmitSourcePrefix bool                  `yaml:"omitSourcePrefix"`
	Global           GlobalCode            `yaml:"global"`
	Groups           []string              `yaml:"groups"`
	Options          map[string]string     `yaml:"options"`
	Keys             map[string]CodeKey    `yaml:"keys"`
	Inject           map[string]CodeInject `yaml:"inject"`
	Locations        map[string]string     `yaml:"exportLocations"`
	FileHeader       string                `yaml:"fileHeader"`
}

type CodeInject struct {
	Name         string   `yaml:"name"`
	Language     string   `yaml:"language"`
	PropertyName string   `yaml:"propertyName"`
	Interface    string   `yaml:"interface"`
	Namespace    string   `yaml:"namespace"`
	Method       string   `yaml:"method"`
	Imports      []string `yaml:"imports"`
}

type GlobalCode struct {
	VariablePrefix  string `yaml:"variablePrefix"`
	VariablePostfix string `yaml:"variablePostfix"`
}
type CodeKey struct {
	// Name      string `yaml:"name"`
	Path      string   `yaml:"path"`
	NameSpace string   `yaml:"namespace"`
	Postfix   string   `yaml:"postfix"`
	Prefix    string   `yaml:"prefix"`
	Imports   []string `yaml:"imports"`
	Inject    []string `yaml:"inject"`
}
