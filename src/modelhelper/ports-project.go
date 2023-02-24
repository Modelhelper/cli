package modelhelper

import (
	"context"
	"modelhelper/cli/modelhelper/models"
)

type ProjectConfigService interface {
	New() (*models.ProjectConfig, error)
	Exists() bool
	Save() error
	Load() (*models.ProjectConfig, error)
	LoadFromFile(path string) (*models.ProjectConfig, error)
	FindReleatedProjects(startPath string) []string
	FindNearestProjectDir() (string, bool)
}

type ProjectGenerator interface {
	Generate(ctx context.Context, template *models.ProjectTemplate, model *models.ProjectTemplateModel) ([]*models.ProjectSourceFile, error)
	BuildTemplateModel(options *models.ProjectTemplateCreateOptions, tpl *models.ProjectTemplate) *models.ProjectTemplateModel
	GenerateRootDirectoryName(rootFolderTemplateName string, model *models.ProjectTemplateModel) (string, error)
}

type ProjectTemplateService interface {
	List(options *models.ProjectTemplateListOptions) map[string]models.ProjectTemplate
	Load(name string) *models.ProjectTemplate
	Group(by string, templateList map[string]models.ProjectTemplate) map[string]map[string]models.ProjectTemplate
}

// type ProjectTemplateModel interface {
// 	models.BasicModel | models.EntityListModel | models.EntityModel
// }
