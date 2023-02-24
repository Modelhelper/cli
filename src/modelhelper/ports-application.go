package modelhelper

import (
	"context"
	"modelhelper/cli/modelhelper/models"
)

type AppInitializer interface {
	IsInitialized() bool
	Initialize() error
}

type AppInfoService interface {
	// LoadConfig() *models.Config
	About() string
	Welcome() string
	Logo() string
	Slogan() string
	Version() string
}

type CommandService interface {
	Execute(ctx context.Context) error
	// BuildCommandTree() error
}

type TemplateFilter interface {
	Filter(t map[string]models.CodeTemplate, filter []string) map[string]models.CodeTemplate
}
