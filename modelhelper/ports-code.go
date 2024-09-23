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

type CodeGeneratorService interface {
	Generate(ctx context.Context, options *models.CodeGeneratorOptions) (*models.CodeGenerateResult, error)
	// GetModel
}

type TemplateTypes interface {
	*models.CodeTemplate | *models.ProjectTemplate | *models.TextTemplate
}

type TemplateGenerator[T TemplateTypes] interface {
	Generate(ctx context.Context, tpl T, mdl interface{}, options *models.CodeTemplateListOptions) (*models.TemplateGeneratorResult, error)
}

type LanguageDefinitionService interface {
	List() map[string]models.LanguageDefinition
	GetDefinition(lang string) *models.LanguageDefinition
}

type CommitHistoryService interface {
	GetCommitHistory(repoPath string, options *models.CommitHistoryOptions) (*models.CommitHistory, error)
	GetTags(repoPath string, options *models.CommitHistoryOptions) ([]models.GitTag, error)
	GetAuthors(repoPath string, options *models.CommitHistoryOptions) (*models.Author, error)
}
