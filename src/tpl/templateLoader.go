package tpl

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type TemplateLoader struct {
	Directory string
}

type TemplateLoaderResult struct {
}

type TemplateMap map[string]Template

// type TemplateDirectory struct {
// 	Directory
// }
func ExtractBlocks(templates *TemplateMap) TemplateMap {
	blocks := make(TemplateMap)

	for k, t := range *templates {

		if (strings.ToLower(t.Type)) == "block" {
			blocks[k] = t
		}
	}

	return blocks
}
func (loader *TemplateLoader) LoadTemplates() (TemplateMap, error) {
	templates := make(TemplateMap)
	path := loader.Directory

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() == false && strings.HasSuffix(p, "yaml") {
			name := convertFileNameToTemplateName(path, p)
			t, err := loadTemplateFromFile(p)
			if err != nil {

			}

			templates[name] = *t
		}

		return nil
	})
	if err != nil {
		log.Println(err)
	}

	return templates, nil
}

// func (l *TemplateLoader) LoadBlock(path string, pattern string) (*[]Template, error) {
// 	return nil, nil
// }

// func (l *TemplateLoader) LoadSnippets(path string, pattern string) (*[]Template, error) {
// 	return nil, nil
// }

// func (l *TemplateLoader) LoadTemplates(path string, pattern string) (*[]Template, error) {
// 	templates := []Template{}
// 	files := loadTemplateFiles(l.Directory, "*")

// 	for k, file := range files {
// 		fmt.Println(k, file)
// 	}

// 	return &templates, nil
// }

// func (l *TemplateLoader) LoadTemplate(name string) (*Template, error) {
// 	files := loadTemplateFiles(l.Directory, "*")

// 	var f = files[name]
// 	// var t Template

// 	if len(f) > 0 {
// 		t, err := loadTemplateFromFile(f)
// 		t.Name = name
// 		if err != nil {
// 			log.Fatalf("cannot load file: %v", err)
// 			return nil, err
// 		}

// 		// t.Name = name
// 		return t, nil
// 	}

// 	return nil, nil
// }

func loadTemplateFromFile(fileName string) (*Template, error) {
	var t *Template

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
