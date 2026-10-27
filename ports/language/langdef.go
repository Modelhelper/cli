package language

import (
	"embed"
	"fmt"
	"github.com/Modelhelper/cli/v3/modelhelper"
	"github.com/Modelhelper/cli/v3/modelhelper/models"
	"io/fs"
	"io/ioutil"
	"log"
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
func (s *langDefService) List() map[string]models.LanguageDefinition {
	return s.definitions()
}

// GetDefinition implements modelhelper.LanguageDefinitionService
func (s *langDefService) GetDefinition(lang string) *models.LanguageDefinition {
	l, ok := s.definitions()[lang]

	if ok {
		return &l
	}

	return nil
}

// definitions returns the definitions built into the CLI with the ones in the folder set as
// languages.definitions in the config laid over them. A language in the folder replaces the
// built in language with the same name, so languages can be added or changed without a rebuild.
func (s *langDefService) definitions() map[string]models.LanguageDefinition {
	defs, _ := loadInternalFiles()

	if s.cfg == nil || len(s.cfg.Languages.Definitions) == 0 {
		return defs
	}

	external, err := loadFromPath(s.cfg.Languages.Definitions)
	if err != nil {
		// the folder does not exist, use the built in definitions
		return defs
	}

	for name, def := range external {
		defs[name] = def
	}

	return defs
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
