package tpl

type TemplateScope string

const (
	GlobalScope  TemplateScope = "global"
	ProjectScope               = "project"
)

type Loader interface {
	LoadTemplates(path string) (*[]Template, error)
}

type TemplateHandler struct {
	Directory string
}
