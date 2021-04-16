package tpl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type Template struct {
	Name            string `yaml:"Name"`
	Version         string `yaml:"Version"`
	InjectKey       string
	Language        string
	LanguageVersion string
	Scope           TemplateScope
	Type            string
	Description     string
	Short           string
	Tags            []string
	Groups          []string
	Import          string
	LocationKey     string
	ExportFileName  string
	Body            string `yaml:"Body"`
}

type TemplateType struct {
	Name      string
	CanExport bool
	IsSnippet bool
}

var fileTemplateType = TemplateType{Name: "file", IsSnippet: true, CanExport: false}

var (
	TemplateTypes = map[string]TemplateType{
		"file":    TemplateType{Name: "file", IsSnippet: false, CanExport: true},
		"snippet": TemplateType{Name: "snippet", IsSnippet: true, CanExport: false},
		"init":    TemplateType{Name: "init", IsSnippet: false, CanExport: false},
	}
)

func (t *Template) Generate(model interface{}) (string, error) {
	blocks := []*Template{
		testBlockLvl1(),
		testBlockLvl2(),
		t,
	}

	t1 := useTempdir(blocks)

	buf := new(bytes.Buffer)
	err := t1.ExecuteTemplate(buf, t.Name, model)

	if err != nil {
		fmt.Println(err)
	}

	return buf.String(), nil
}

func withoutTempDir(name string, blocks []*Template) *template.Template {

	t := template.New(name)

	for _, b := range blocks {
		t = template.Must(t.Parse(b.Body))
	}
	return t
}

func useTempdir(blocks []*Template) *template.Template {

	dir := createTempDir()
	defer os.RemoveAll(dir)

	err := writeTempFiles(dir, blocks)

	if err != nil {
		return nil
	}

	pattern := filepath.Join(dir, "*")

	drivers := template.Must(template.ParseGlob(pattern))

	return drivers
}

func createTempDir() string {

	//err := os.Mkdir("t", 0777)
	fmt.Println(os.TempDir())

	dir, err := ioutil.TempDir("", "template") //os.MkdirTemp("", "template")
	fmt.Println(dir)
	if err != nil {
		log.Fatal(err)
	}

	return dir

	// for _, file := range files {
	// 	fp := filepath.Join(dir, file.Name)
	// 	fmt.Println(fp)
	// 	err := ioutil.WriteFile(fp, []byte(file.Body), 0777)
	// 	// f, err := os.Create(fp)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	// defer f.Close()+".tmpl"
	// 	// _, err = io.WriteString(f, file.Body)
	// 	// if err != nil {
	// 	// 	log.Fatal(err)
	// 	// }
	// }
	// return dir
}

func writeTempFiles(dir string, files []*Template) error {
	for _, file := range files {
		fp := filepath.Join(dir, file.Name)
		err := ioutil.WriteFile(fp, []byte(file.Body), 0777)
		if err != nil {
			return err
		}
	}

	return nil
}
