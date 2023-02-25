package language

import (
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed defs
var langDefFs embed.FS

type langDefService struct {
	cfg *models.Config
}

// List implements modelhelper.LanguageDefinitionService
func (*langDefService) List() map[string]models.LanguageDefinition {
	files, _ := loadInternalFiles()
	return files
}

// GetDefinition implements modelhelper.LanguageDefinitionService
func (*langDefService) GetDefinition(lang string) *models.LanguageDefinition {
	langs, _ := loadInternalFiles()

	l, ok := langs[lang]

	if ok {
		return &l
	}

	return nil
}

func NewLanguageDefinitionService(cfg *models.Config) modelhelper.LanguageDefinitionService {
	return &langDefService{cfg}
}

func loadInternalFiles() (map[string]models.LanguageDefinition, error) {

	defs := make(map[string]models.LanguageDefinition)

	files, _ := langDefFs.ReadDir("defs")

	for _, ff := range files {

		var langDef *models.LanguageDefinition
		fname := fmt.Sprintf("defs/%s", ff.Name())
		bytes, _ := fs.ReadFile(langDefFs, fname) // langDefFs.ReadFile(ff.Name())
		err := yaml.Unmarshal(bytes, &langDef)
		if err != nil {
			log.Fatalf("cannot unmarshal data: %v", err)
			// return nil, err
		}

		if langDef != nil {
			defs[langDef.Language] = *langDef
		}

	}

	return defs, nil

}
func loadFromPath(dir string) (map[string]models.LanguageDefinition, error) {

	defs := make(map[string]models.LanguageDefinition)

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

func loadDef(fileName string) (*models.LanguageDefinition, error) {
	var langDef *models.LanguageDefinition
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
