package modelhelper

import (
	"modelhelper/cli/modelhelper/models"
)

type CreateConfigService interface {
	Create() *models.Config
}
type ConfigService interface {
	Load() (*models.Config, error)
	Save() error
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
