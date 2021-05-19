package project

import (
	"io/ioutil"
	"log"
	"modelhelper/cli/source"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const dirname string = ".modelhelper"

type Project struct {
	Version       string                       `yaml:"version"`
	Name          string                       `yaml:"name"`
	Language      string                       `yaml:"language"`
	DefaultSource string                       `yaml:"defaultSource,omitempty"`
	DefaultKey    string                       `yaml:"defaultKey,omitempty"`
	Connections   map[string]source.Connection `yaml:"connections,omitempty"`
	Code          ProjectCode                  `yaml:"code,omitempty"`
	OwnerName     string                       `yaml:"customerName,omitempty"`
	Header        string                       `yaml:"header,omitempty"`
	Options       map[string]string            `yaml:"options,omitempty"`
	Custom        interface{}                  `yaml:"custom,omitempty"`
}

func (p *Project) Save() error {

	d, err := yaml.Marshal(&p)

	if err != nil {

		return err
	}

	// path := filepath.Join(DefaultLocation(), "config.yaml")
	path := DefaultLocation()

	if !Exists(path) {
		CreateDir(dirname)
	}

	err = ioutil.WriteFile(path, d, 0777)

	return err
}

func CreateDir(name string) {
	err := os.Mkdir(name, 0755)
	if err != nil {
		log.Fatal(err)
	}
}
func DefaultDir() string {
	p, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return filepath.Join(p, dirname)
}
func DefaultLocation() string {
	// p, err := os.Getwd()
	// if err != nil {
	// 	log.Println(err)
	// }
	return filepath.Join(DefaultDir(), "project.yaml")
}

func (P *Project) Exists(path string) bool {

	pathInfo, err := os.Stat(path)

	if os.IsNotExist(err) || pathInfo.IsDir() {
		return false
	}

	return true
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
			// log.Fatal("Project does not exits")
			return nil, err
		}

	} else {
		p, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		path = filepath.Join(p, dirname, "project.yaml")
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
	Keys             map[string]CodeKey    `yaml:"keys,omitempty"`
	Inject           map[string]CodeInject `yaml:"inject,omitempty"`
	Locations        map[string]string     `yaml:"exportLocations"`
	FileHeader       string                `yaml:"fileHeader"`
}

type CodeInject struct {
	Name         string   `yaml:"name,omitempty"`
	Language     string   `yaml:"language,omitempty"`
	PropertyName string   `yaml:"propertyName,omitempty"`
	Interface    string   `yaml:"interface,omitempty"`
	Namespace    string   `yaml:"namespace,omitempty"`
	Method       string   `yaml:"method,omitempty"`
	Imports      []string `yaml:"imports,omitempty"`
}

type GlobalCode struct {
	VariablePrefix  string `yaml:"variablePrefix"`
	VariablePostfix string `yaml:"variablePostfix"`
}
type CodeKey struct {
	// Name      string `yaml:"name"`
	Path      string   `yaml:"path,omitempty"`
	NameSpace string   `yaml:"namespace,omitempty"`
	Postfix   string   `yaml:"postfix,omitempty"`
	Prefix    string   `yaml:"prefix,omitempty"`
	Imports   []string `yaml:"imports,omitempty"`
	Inject    []string `yaml:"inject,omitempty"`
}

func FindNearestProjectDir() (string, bool) {
	basePath := "./"
	root := DefaultDir()
	updir := "../"

	for i := 1; i < 6; i++ {
		files, err := ioutil.ReadDir(basePath)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {

			if f.IsDir() && f.Name() == dirname {
				fp := filepath.Join(basePath, f.Name(), "project.yaml")
				return fp, true
			}
		}

		relp, _ := filepath.Rel(root, basePath)
		basePath = filepath.Join(updir, relp)

		updir += "../"
	}

	return "", false
}
