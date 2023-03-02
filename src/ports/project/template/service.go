package template

import (
	"fmt"
	"io/ioutil"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/constants"
	"modelhelper/cli/modelhelper/models"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type projectTemplateService struct {
	config *models.Config
}

// List implements modelhelper.ProjectTemplateService
func (p *projectTemplateService) List(options *models.ProjectTemplateListOptions) map[string]models.ProjectTemplate {

	files := p.getProjectFileList()
	return createTemplateMap(files)
}

// Group implements modelhelper.ProjectTemplateService
func (p *projectTemplateService) Group(by string, templateList map[string]models.ProjectTemplate) map[string]map[string]models.ProjectTemplate {
	panic("unimplemented")
}

func (p *projectTemplateService) getProjectFileList() []string {
	ls := []string{}

	for _, location := range p.config.Templates.Project {
		files := getProjectTemplateFiles(location)
		ls = append(ls, files...)

	}

	return ls
}

// Load implements modelhelper.ProjectTemplateService
func (p *projectTemplateService) Load(name string) *models.ProjectTemplate {
	files := p.getProjectFileList()

	for _, file := range files {
		pt, err := loadTemplateFromFile(file)
		if err != nil {
			fmt.Printf("\nCould not load file '%s'\n\t%v", file, err)
		}

		if pt != nil && len(pt.Name) > 0 && strings.ToLower(pt.Name) == strings.ToLower(name) {
			return pt
		}
	}

	return nil

}

// Load implements modelhelper.ProjectTemplateService
func createTemplateMap(files []string) map[string]models.ProjectTemplate {
	m := make(map[string]models.ProjectTemplate)
	for _, file := range files {
		pt, err := loadTemplateFromFile(file)
		if err != nil {
			fmt.Printf("\nCould not load file '%s'\n\t%v", file, err)
		}

		if pt != nil && len(pt.Name) > 0 {
			m[pt.Name] = *pt
		}
	}
	return m
}

func loadTemplateFromFile(fileName string) (*models.ProjectTemplate, error) {
	var t *models.ProjectTemplate

	dat, e := ioutil.ReadFile(fileName)
	if e != nil {
		return nil, e
	}

	err := yaml.Unmarshal(dat, &t)
	if err != nil {
		return nil, err
	}

	if t != nil {
		tDir, fName := filepath.Split(fileName)
		t.TemplateFilePath = tDir
		t.TemplateFileName = fName
	}
	return t, nil
}

func NewProjectTemplateService(cfg *models.Config) modelhelper.ProjectTemplateService {
	return &projectTemplateService{cfg}
}

func getProjectTemplateFiles(path string) []string {
	files := []string{}

	filepath.Walk(path, func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(fullPath, fmt.Sprintf("%s.yaml", constants.ProjectTemplateFileName)) ||
			strings.HasSuffix(fullPath, fmt.Sprintf("%s.yml", constants.ProjectTemplateFileName))) {
			files = append(files, fullPath)
		}

		return nil
	})

	return files
}
