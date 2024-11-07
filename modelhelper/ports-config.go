package modelhelper

import (
	"modelhelper/cli/modelhelper/models"
)

type CreateConfigService interface {
	Create() *models.Config
}
type ConfigService interface {
	ConfigExists() bool
	Create() (*models.Config, error)
	Load() (*models.Config, error)
	Save() error
	SaveConfig(c *models.Config) error
	LoadFromFile(path string) (*models.Config, error)
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
