package generator

import (
	"context"
	"modelhelper/cli/modelhelper"
)

type CodeGeneratorOptions struct {
	Templates         []string
	TemplateGroups    []string
	TemplatePath      string
	CanUseTemplates   bool
	EntityGroups      []string
	Entities          []string
	ExportToScreen    bool
	ExportByKey       bool
	ExportPath        string
	Connection        string
	ExportToClipboard bool
	Overwrite         bool
	Relations         string
	CodeOnly          bool
	UseDemo           bool
	ConfigFilePath    string
	ProjectFilePath   string
}

type codeGenerator struct {
	options         *CodeGeneratorOptions
	templateService modelhelper.CodeTemplateService
}

func NewCodeGenerator(options *CodeGeneratorOptions, templateService modelhelper.CodeTemplateService) modelhelper.CodeGenerator {
	return &codeGenerator{
		options:         options,
		templateService: templateService,
	}
}

func (g *codeGenerator) Generate(ctx context.Context) ([]modelhelper.CodeFileResult, error) {

	return nil, nil
}
