package modelhelper

import (
	"context"
	"modelhelper/cli/modelhelper/models"
)

type ConfigLoader interface {
	Load() (*models.Config, error)
	LoadFromFile(path string) (*models.Config, error)
	Save() error
	// GetConnections() (*map[string]Connection, error)
}

type ProjectService interface {
	New() (*models.ProjectConfig, error)
	Exists() bool
	Save() error
	Load() (*models.ProjectConfig, error)
	LoadFromFile(path string) (*models.ProjectConfig, error)
	FindReleatedProjects(startPath string) []string
	FindNearestProjectDir() (string, bool)
}

type CodeTemplateService interface {
	List() map[string]models.CodeTemplate
	Load(name string) *models.CodeTemplate
}

type CodeModel interface {
	models.BasicModel | models.EntityListModel | models.EntityModel
}

type CodeGenerator interface {
	Generate(ctx context.Context) ([]models.CodeFileResult, error)
}

type SourceService interface {
	ConnectionStringPart(part string) string
	ParseConnectionString()
	Entity(name string) (*models.Entity, error)
	Entities(pattern string) (*[]models.Entity, error)
	EntitiesFromColumn(column string) (*[]models.Entity, error)
}

type ConnectionProvider interface {
	GetConnections() (*map[string]models.Connection, error)
}

type AppService interface {
	LoadConfig() *models.Config
	About() string
	Logo() string
	Slogan() string
}

// type Services struct {
// 	ConfigLoader ConfigLoader
// 	Project      ProjectService
// 	Code         CodeGenerator
// }

// func NewModelhelperApp() *App {
// 	return &App{
// 		ConfigLoader: config.NewConfigLoader(),
// 		Project:      project.NewModelhelperProject(),
// 	}
// }
