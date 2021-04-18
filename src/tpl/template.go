package tpl

type TemplateScope string

const (
	GlobalScope  TemplateScope = "global"
	ProjectScope               = "project"
)

// obsolete
type LoadHandler interface {
	LoadBlock(path string, pattern string) (*[]Template, error)
	LoadSnippets(path string, pattern string) (*[]Template, error)
	LoadTemplates(path string, pattern string) (*[]Template, error)
	LoadTemplate(name string) (*Template, error)
}

// obsolete
type Generator interface {
	Generate(model interface{}) (string, error)
}

type TemplateHandler struct {
	Directory string
}
