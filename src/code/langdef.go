package code

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type LanguageDefinition struct {
	Version        string              `json:"version" yaml:"version"`
	Language       string              `json:"language" yaml:"language"`
	DataTypes      map[string]Datatype `json:"datatypes" yaml:"datatypes"`
	DefaultImports []string            `json:"defaultImports" yaml:"defaultImports"`
	Keys           map[string]Key      `json:"keys" yaml:"keys"`
	Inject         map[string]Inject   `json:"inject" yaml:"inject"`
	Global         Global              `json:"global" yaml:"global"`
	Short          string              `json:"short" yaml:"short"`
	Description    string              `json:"description" yaml:"description"`
	Path           string
	// CanInject                 bool                       `json:"canInject" yaml:"canInject"`
	// UsesNamespace             bool                       `json:"usesNamespace" yaml:"usesNamespace"`
	// ModuleLevelVariablePrefix string                     `json:"moduleLevelVariablePrefix" yaml:"moduleLevelVariablePrefix"`
}

type Code struct {
	RootNamespace          string            `yaml:"rootNamespace,omitempty"`
	OmitSourcePrefix       bool              `yaml:"omitSourcePrefix,omitempty"`
	Global                 Global            `yaml:"global"`
	Groups                 []string          `yaml:"groups"`
	Options                map[string]string `yaml:"options"`
	Keys                   map[string]Key    `yaml:"keys,omitempty"`
	Inject                 map[string]Inject `yaml:"inject,omitempty"`
	Locations              map[string]string `yaml:"locations"`
	FileHeader             string            `yaml:"header"`
	DisableNullableTypes   bool              `json:"diableNullableTypes" yaml:"diableNullableTypes"`
	UseNullableAlternative bool              `json:"useNullableAlternative" yaml:"useNullableAlternative"`
}
type Datatype struct {
	Key                 string      `json:"key" yaml:"key"`
	NotNull             string      `json:"notNull" yaml:"notNull"`
	Nullable            string      `json:"nullable" yaml:"nullable"`
	NullableAlternative string      `json:"nullableAlternative" yaml:"nullableAlternative"`
	DefaultValue        interface{} `json:"defaultValue" yaml:"defaultValue"`
}

type Inject struct {
	Name         string   `json:"name" yaml:"name"`
	PropertyName string   `json:"propertyName" yaml:"propertyName"`
	Method       string   `json:"method" yaml:"method"`
	Imports      []string `json:"imports" yaml:"imports"`
}

type Key struct {
	Postfix   string   `json:"postfix" yaml:"postfix"`
	Prefix    string   `json:"prefix" yaml:"prefix"`
	Imports   []string `json:"imports" yaml:"imports"`
	Inject    []string `json:"inject" yaml:"inject"`
	Namespace string   `json:"namespace" yaml:"namespace"`
}

type Global struct {
	VariablePrefix  string `yaml:"variablePrefix"`
	VariablePostfix string `yaml:"variablePostfix"`
}

// func Load() (map[string]LanguageDefinition, error) {
// 	cfg := config.Load()

// 	if len(cfg.Languages.Definitions) == 0 {
// 		return nil, nil
// 	}
// 	return LoadFromPath(cfg.Languages.Definitions)
// }

func LoadFromPath(dir string) (map[string]LanguageDefinition, error) {

	defs := make(map[string]LanguageDefinition)

	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		return nil, err
	}

	files, _ := ioutil.ReadDir(dir)

	for _, file := range files {
		fileName := filepath.Join(dir, file.Name())
		if !file.IsDir() && (strings.HasSuffix(fileName, "yaml") || strings.HasSuffix(fileName, "yml")) {
			t, err := loadDef(fileName)

			if err != nil {
				log.Fatalln(err)
			}

			if t != nil {
				defs[t.Language] = *t
			}
		}
	}

	return defs, nil

}

func loadDef(fileName string) (*LanguageDefinition, error) {
	var langDef *LanguageDefinition
	// fmt.Println("load file: ", fileName)
	dat, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Fatalf("cannot load file: %v", e)
		return nil, e
	}

	err := yaml.Unmarshal(dat, &langDef)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
		return nil, err
	}

	if langDef == nil {
		return nil, nil
	}

	langDef.Path = fileName
	return langDef, nil
}
