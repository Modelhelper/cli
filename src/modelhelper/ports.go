package modelhelper

import (
	"context"
)

type ConfigLoader interface {
	Load() (*Config, error)
	LoadFromFile(path string) (*Config, error)
	Save() error
	// GetConnections() (*map[string]Connection, error)
}

type ProjectService interface {
	New() (*ProjectConfig, error)
	Exists() bool
	Save() error
	Load() (*ProjectConfig, error)
	LoadFromFile(path string) (*ProjectConfig, error)
	FindReleatedProjects(startPath string) []string
	FindNearestProjectDir() (string, bool)
}

type CodeTemplateService interface {
	List() map[string]CodeTemplate
	Load(name string) *CodeTemplate
}

type CodeModel interface {
	BasicModel | EntityListModel | EntityModel
}

type CodeGenerator interface {
	Generate(ctx context.Context) ([]CodeFileResult, error)
}

type SourceService interface {
	ConnectionStringPart(part string) string
	ParseConnectionString()
	Entity(name string) (*Entity, error)
	Entities(pattern string) (*[]Entity, error)
	EntitiesFromColumn(column string) (*[]Entity, error)
}

type ConnectionProvider interface {
	GetConnections() (*map[string]Connection, error)
}

type AppService interface {
	LoadConfig() *Config
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
