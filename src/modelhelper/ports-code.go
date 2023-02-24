package modelhelper

import (
	"context"
	"modelhelper/cli/modelhelper/models"
)

type CodeTemplateService interface {
	List(options *models.CodeTemplateListOptions) map[string]models.CodeTemplate
	Load(name string) *models.CodeTemplate
	Group(by string, templateList map[string]models.CodeTemplate) map[string]map[string]models.CodeTemplate
}

type CodeModel interface {
	models.BasicModel | models.EntityListModel | models.EntityModel
}

type CodeGenerator interface {
	Generate(ctx context.Context, options *models.CodeGeneratorOptions) ([]models.TemplateGeneratorFileResult, error)
	GenerateCode(tpl *models.CodeTemplate, mdl interface{}) (*models.TemplateGeneratorFileResult, error)
}

type TemplateTypes interface {
	*models.CodeTemplate | *models.ProjectTemplate | *models.TextTemplate
}

type TemplateGenerator[T TemplateTypes] interface {
	Generate(ctx context.Context, tpl T, mdl interface{}) (*models.TemplateGeneratorResult, error)
}
