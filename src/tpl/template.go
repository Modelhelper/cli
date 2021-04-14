package tpl

import (
	"bytes"
	"text/template"
)

type TemplateScope string

const (
	GlobalScope  TemplateScope = "global"
	ProjectScope               = "project"
)

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

type Generator interface {
	Generate(model interface{}) (string, error)
}

func (t *Template) Generate(model interface{}) (string, error) {

	tpl := template.Must(template.New(t.Language).Parse(t.Body))
	buf := new(bytes.Buffer)

	// buf.String() // returns a string of what was written to it
	tpl.Execute(buf, model)

	return buf.String(), nil
}
