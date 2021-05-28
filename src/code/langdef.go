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
	Version        string                     `json:"version" yaml:"version"`
	Language       string                     `json:"language" yaml:"language"`
	DataTypes      map[string]LangDefDataType `json:"datatypes" yaml:"datatypes"`
	DefaultImports []string                   `json:"defaultImports" yaml:"defaultImports"`
	// CanInject                 bool                       `json:"canInject" yaml:"canInject"`
	// UsesNamespace             bool                       `json:"usesNamespace" yaml:"usesNamespace"`
	// ModuleLevelVariablePrefix string                     `json:"moduleLevelVariablePrefix" yaml:"moduleLevelVariablePrefix"`
}

type LangDefDataType struct {
	Key                 string `json:"key" yaml:"key"`
	NotNull             string `json:"notNull" yaml:"notNull"`
	Nullable            string `json:"nullable" yaml:"nullable"`
	NullableAlternative string `json:"nullableAlternative" yaml:"nullableAlternative"`
}

type LangDefInject struct {
	Name         string   `json:"name" yaml:"name"`
	PropertyName string   `json:"propertyName" yaml:"propertyName"`
	Imports      []string `json:"imports" yaml:"imports"`
}

type LangDefKey struct {
	Postfix   string   `json:"postfix" yaml:"postfix"`
	Prefix    string   `json:"prefix" yaml:"prefix"`
	Imports   []string `json:"imports" yaml:"imports"`
	Inject    []string `json:"inject" yaml:"inject"`
	Namespace string   `json:"namespace" yaml:"namespace"`
}

func LoadFromPath(dir string) (map[string]LanguageDefinition, error) {

	defs := make(map[string]LanguageDefinition)

	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() == false && strings.HasSuffix(p, "yaml") {
			t, err := loadDef(p)

			if err != nil {
				log.Fatalln(err)
			}

			defs[t.Language] = *t
		}

		return nil
	})
	if err != nil {
		log.Println(err)
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

	return langDef, nil
}
