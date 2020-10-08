package types

type TemplateScope string

const (
	GlobalScope  TemplateScope = "global"
	ProjectScope               = "project"
)

type Template struct {
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
