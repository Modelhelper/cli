package app

import "modelhelper/cli/config"

func Initialize(initializer Initializer) error {
	err := initializer.Initialize()
	return err
}

type Initializer interface {
	Initialize() error
}

func (a *Application) Initialize() error {

	PrintWelcomeMessage()

	cfg := config.New()

	if err := cfg.Initialize(); err != nil {
		return err
	}

	return nil
}
