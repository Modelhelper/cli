package template

import (
	"io"
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

type codeTemplateService struct {
	config *models.Config
}

type codeFile struct {
	fullPath        string
	fileNameFromDir string
}

func NewCodeTemplateService(cfg *models.Config) modelhelper.CodeTemplateService {
	return &codeTemplateService{cfg}
}

// List implements modelhelper.CodeTemplateService
func (t *codeTemplateService) List(options *models.CodeTemplateListOptions, append ...fs.FS) map[string]models.CodeTemplate {
	if len(t.config.Templates.Code) == 0 {
		return nil
	}

	templates := make(map[string]models.CodeTemplate)
	// path := t.config.Templates.Code

	for _, codeFile := range t.getFileList(options) {
		// name := convertFileNameToTemplateName(path, p)
		t, err := loadTemplateFromFile(codeFile.fullPath)
		if err != nil {
			log.Fatalln(err)
		}

		if t != nil {
			if codeFile != nil {
				t.Name = codeFile.fileNameFromDir
			}
			templates[codeFile.fileNameFromDir] = *t
		}
	}

	if len(append) > 0 {
		for _, a := range append {
			ts := getCodeTemplatesFromFS(a)

			for k, t := range ts {
				templates[k] = t
			}
		}
	}

	if options != nil {
		if len(options.FilterKeys) > 0 {
			templates = filter("key", templates, options.FilterKeys)
		}
		if len(options.FilterLanguages) > 0 {
			templates = filter("language", templates, options.FilterLanguages)
		}
		if len(options.FilterTypes) > 0 {
			templates = filter("type", templates, options.FilterTypes)
		}
		if len(options.FilterGroups) > 0 {
			templates = filter("groups", templates, options.FilterGroups)
		}
		if len(options.FilterModels) > 0 {
			templates = filter("model", templates, options.FilterModels)
		}

	}

	return templates
}

func loadTemplateFromBytes(data []byte, isEmbedded bool) (*models.CodeTemplate, error) {
	var t *models.CodeTemplate

	err := yaml.Unmarshal(data, &t)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
		return nil, err
	}

	if t != nil {
		t.TemplateFilePath = ""
		t.IsEmbedded = isEmbedded
	}
	return t, nil
}

func (p *codeTemplateService) getFileList(options *models.CodeTemplateListOptions) []*codeFile {
	ls := []*codeFile{}

	for _, location := range p.config.Templates.Code {
		files := getCodeTemplateFiles(location)
		ls = append(ls, files...)

	}

	if p.config.Templates.Database != nil && options != nil && options.DatabaseType != "" {

		for _, dbPath := range p.config.Templates.Database {
			files := getDatabaseTemplateFiles(dbPath, options.DatabaseType)
			ls = append(ls, files...)
		}
	}

	return ls
}

func getCodeTemplateFiles(path string) []*codeFile {
	files := []*codeFile{}

	filepath.Walk(path, func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(fullPath, "yaml") || strings.HasSuffix(fullPath, "yml")) {

			cf := &codeFile{
				fullPath:        fullPath,
				fileNameFromDir: convertFileNameToTemplateName(path, fullPath),
			}
			files = append(files, cf)
		}

		return nil
	})

	return files
}

func getCodeTemplatesFromFS(fileSys fs.FS) map[string]models.CodeTemplate {
	templates := make(map[string]models.CodeTemplate)

	fs.WalkDir(fileSys, ".", func(path string, d fs.DirEntry, err error) error {

		if !d.IsDir() && (strings.HasSuffix(d.Name(), "yaml") || strings.HasSuffix(d.Name(), "yml")) {
			name := convertFileNameToTemplateName(path, d.Name())

			file, err := fileSys.Open(d.Name())
			if err != nil {
				log.Fatalln(err)
			}

			data, err := io.ReadAll(file)
			if err != nil {
				log.Fatalln(err)
			}

			tpl, err := loadTemplateFromBytes(data, true)

			if err != nil {
				log.Fatalln(err)
			}

			templates[name] = *tpl

		}

		return nil
	})

	// for _, codeFile := range files {
	// 	// name := convertFileNameToTemplateName(path, p)

	// 	t, err := loadTemplateFromFile(codeFile.fullPath)
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}

	// 	if t != nil {
	// 		if codeFile != nil {
	// 			t.Name = codeFile.fileNameFromDir
	// 		}
	// 		templates[codeFile.fileNameFromDir] = *t
	// 	}
	// }
	return templates
}

func getDatabaseTemplateFiles(path, dbType string) []*codeFile {
	typeConverter := make(map[string]string)
	typeConverter["sqlserver"] = "mssql"
	typeConverter["mssql"] = "mssql"
	typeConverter["mysql"] = "mysql"
	typeConverter["postgres"] = "postgres"
	typeConverter["postgresql"] = "postgres"
	typeConverter["pg"] = "postgres"

	dbType = typeConverter[dbType]
	path = filepath.Join(path, dbType)
	return getCodeTemplateFiles(path)
}

func filter(filterType string, t map[string]models.CodeTemplate, filter []string) map[string]models.CodeTemplate {
	output := make(map[string]models.CodeTemplate)

	for name, template := range t {
		switch filterType {
		case "type":
			if contains(filter, template.Type) {
				output[name] = template
			}
		case "key":
			if contains(filter, template.Key) {
				output[name] = template
			}
		case "model":
			if contains(filter, template.Model) {
				output[name] = template
			}
		case "language":
			if contains(filter, template.Language) {
				output[name] = template
			}
		case "groups":
			if len(template.Features) > 0 {
				for _, grp := range template.Features {

					if contains(filter, grp) {
						output[name] = template
					}
				}
			}
		}
	}

	return output
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (cts *codeTemplateService) Group(by string, templateList map[string]models.CodeTemplate) map[string]map[string]models.CodeTemplate {
	m := make(map[string]map[string]models.CodeTemplate)
	empty := make(map[string]models.CodeTemplate)

	for tn, template := range templateList {

		var tmpl = template
		key := ""

		if len(by) > 0 {
			switch by {
			case "language":
				key = template.Language
			// case "group":
			// 	key = template.Groups[]
			case "key":
				key = template.Key
			case "model":
				key = template.Model
			// case "tag":
			// 	return &GroupByTag{}
			default:
				key = template.Type

			}
		}

		if len(key) == 0 {
			empty[tn] = tmpl
		} else {
			k, f := m[key]

			if !f {
				k = make(map[string]models.CodeTemplate)
			}

			k[tn] = template
			m[key] = k
		}
	}

	if len(empty) > 0 {
		m["empty"] = empty
	}

	return m
}

func grouper[T []string](list T, templateList map[string]models.CodeTemplate) map[string]map[string]models.CodeTemplate {

	m := make(map[string]map[string]models.CodeTemplate)
	empty := make(map[string]models.CodeTemplate)

	for tn, template := range templateList {

		if template.Type != "block" {

			var tmpl = template

			if len(list) == 0 {
				empty[tn] = tmpl
			} else {
				for _, grp := range list {

					k, f := m[grp]

					if !f {
						k = make(map[string]models.CodeTemplate)
					}

					k[tn] = tmpl
					m[grp] = k
				}
			}
		}
	}

	if len(empty) > 0 {
		m["empty"] = empty
	}

	return m
}

func listGrouper(key string, template models.CodeTemplate, list []string) map[string]map[string]models.CodeTemplate {
	m := make(map[string]map[string]models.CodeTemplate)
	empty := make(map[string]models.CodeTemplate)

	if len(list) == 0 {
		empty[key] = template
	} else {
		for _, grp := range list {

			k, f := m[grp]

			if !f {
				k = make(map[string]models.CodeTemplate)
			}

			k[key] = template
			m[grp] = k
		}
	}

	if len(empty) > 0 {
		m["empty"] = empty
	}

	return m
}

// Load implements modelhelper.CodeTemplateService
func (t *codeTemplateService) Load(name string) *models.CodeTemplate {

	return nil
}

func loadTemplateFiles(path string, pattern string) map[string]string {
	fa := make(map[string]string)

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() == false && strings.HasSuffix(p, "yaml") {
			name := convertFileNameToTemplateName(path, p)
			fa[name] = p
			// fa = append(fa, p)
		}

		// fmt.Println(p, info.Size(), info.Name(), info.IsDir(), info.Sys())  && info.Name()[:3] == "txt"
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	return fa
}

func loadTemplateFromFile(fileName string) (*models.CodeTemplate, error) {
	var t *models.CodeTemplate

	dat, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Fatalf("cannot load file: %v", e)
		return nil, e
	}

	err := yaml.Unmarshal(dat, &t)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
		return nil, err
	}

	if t != nil {
		t.TemplateFilePath = fileName
		t.IsEmbedded = false
	}
	return t, nil
}

func convertFileNameToTemplateName(rootPath string, name string) string {

	out := strings.Replace(name, rootPath, "", -1)
	out = strings.Replace(out, "\\", "-", -1)
	out = strings.Replace(out, "/", "-", -1)
	out = strings.Replace(out, " ", "-", -1)

	var start = 0
	var end = len(out)

	if out[0] == '-' {
		start = 1
	}

	i := strings.LastIndex(out, ".")
	if i > -1 {
		end = i
	}

	return strings.ToLower(out[start:end])
}
