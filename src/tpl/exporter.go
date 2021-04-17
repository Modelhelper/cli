package tpl

type TemplateExporter interface {
	Export() (interface{}, error)
}
