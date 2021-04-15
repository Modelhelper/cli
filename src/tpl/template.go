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

type TemplateScope string

const (
	GlobalScope  TemplateScope = "global"
	ProjectScope               = "project"
)

type Handler interface {
	GetBlock()
	GetSnippets()
	GetTemplates()
}

type Generator interface {
	Generate(model interface{}) (string, error)
}

type TemplateHandler struct {
}
type Template struct {
	Name            string
	Version         string
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
	Body            string
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

	dir := createTestDir(blocks)
	// defer os.RemoveAll(dir)

	pattern := filepath.Join(dir, "*")

	drivers := template.Must(template.ParseGlob(pattern))

	buf := new(bytes.Buffer)
	err := drivers.ExecuteTemplate(buf, t.Name, model)

	// blockTemplates := []*template.Template{} +".tmpl"

	// for _, block := range blocks {
	// 	blockTemplate, err := template.New(block.Name).Parse(block.Body)
	// 	if err != nil {
	// 		continue
	// 	}

	// 	blockTemplates = append(blockTemplates, blockTemplate)
	// }

	// tpl, err := template.New(t.Language).Parse(t.Body)
	if err != nil {
		fmt.Println(err)
	}
	//tpl := template.Must(template.New(t.Language).Parse(t.Body))

	// buf.String() // returns a string of what was written to it
	// tpl.Execute(buf, model)

	return buf.String(), nil
}

func oneMoreTest() *Template {
	t := Template{}

	t.Language = "poco"
	t.Body = testTemplate()

	return &t
}

func testBlockLvl1() *Template {

	t := Template{}

	t.Name = "classname"
	t.Body = `{{ .Name }}`

	return &t
}

func testBlockLvl2() *Template {

	t := Template{}

	t.Name = "namespace"
	t.Body = `namespace {{ .NameSpace }}`

	return &t
}

func createTestDir(files []*Template) string {

	//err := os.Mkdir("t", 0777)
	fmt.Println(os.TempDir())

	dir, err := ioutil.TempDir("", "template") //os.MkdirTemp("", "template")
	fmt.Println(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fp := filepath.Join(dir, file.Name)
		fmt.Println(fp)
		err := ioutil.WriteFile(fp, []byte(file.Body), 0777)
		// f, err := os.Create(fp)
		if err != nil {
			log.Fatal(err)
		}
		// defer f.Close()+".tmpl"
		// _, err = io.WriteString(f, file.Body)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}
	return dir
}

func testTemplate() string {
	return `
// code her - indent by SPACE not TAB
{{- $model := index .Code.Types "model" }}
using System;
{{- range .Code.Imports }}
{{.}}
{{- end }}


//namespace from template
namespace {{ $model.NameSpace }} 
{	
	{{ template "block0" . }}
	public class {{ .Name }}{{ $model.NamePostfix}} 
	{

	}
}

{{.}}
`
}
