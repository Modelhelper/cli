package modelhelper

import "modelhelper/cli/modelhelper/models"

type SourceService interface {
	ConnectionStringPart(part string) string
	ParseConnectionString()
	Entity(name string) (*models.Entity, error)
	Entities(pattern string) (*[]models.Entity, error)
	EntitiesFromNames(names []string) (*[]models.Entity, error)
	EntitiesFromColumn(column string) (*[]models.Entity, error)
}

type ConnectionProvider interface {
	GetConnections() (*map[string]models.Connection, error)
}

type ConnectionBuilder interface {
	Build() *models.Connection
}

type CodeModelConverter interface {
	ToBasicModel(identifier, language string, project *models.ProjectConfig) *models.BasicModel
	ToEntityModel(key, language string, project *models.ProjectConfig, entity *models.Entity) *models.EntityModel
	ToEntityListModel(key, language string, project *models.ProjectConfig, entity *[]models.Entity) *models.EntityListModel
}

type ProjectModelConverter interface {
	ToProjectModel(cfg *models.Config, options *models.ProjectTemplateCreateOptions) *models.ProjectTemplateModel
}
