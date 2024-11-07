package app

import (
	"context"
	"modelhelper/cli/modelhelper"

	"github.com/charmbracelet/huh"
)

type appInitializer struct {
	info        modelhelper.AppInfoService
	cfg         modelhelper.ConfigService
	connService modelhelper.ConnectionService
}

func NewCliInitializer(ctx context.Context, info modelhelper.AppInfoService, cfg modelhelper.ConfigService, cs modelhelper.ConnectionService) modelhelper.AppInitializer {
	return &appInitializer{info, cfg, cs}
}

func (ai *appInitializer) IsInitialized() bool {
	return true
}

func (ai *appInitializer) Initialize() error {
	var (
		codeEditor string
	)
	dbTypeSelect := huh.NewSelect[string]().
		Title("It looks like this is the first time the ModelHelper is run on this computer").
		Description("Please select the type of editor you want to connect use to view things in").
		Options(
			huh.NewOption("I do not want to set a default one", "none"),
			huh.NewOption("Visual Studio Code (code)", "code"),
			huh.NewOption("Vim", "vim"),
			huh.NewOption("Nano", "nano"),
		).
		Value(&codeEditor)

	huh.NewForm(
		huh.NewGroup(dbTypeSelect),
		// huh.NewGroup(ringSelect),
	).Run()

	return nil
}
