package app

import (
	"context"
	"modelhelper/cli/modelhelper"
)

type appInitializer struct {
	info modelhelper.AppInfoService
	cfg  modelhelper.ConfigService
}

func NewCliInitializer(ctx context.Context, info modelhelper.AppInfoService, cfg modelhelper.ConfigService) modelhelper.AppInitializer {
	return &appInitializer{info, cfg}
}

func (ai *appInitializer) IsInitialized() bool {
	return true
}

func (ai *appInitializer) Initialize() error {

	return nil
}
