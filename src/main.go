package main

import (
	"context"
	"fmt"
	"modelhelper/cli/modelhelper"
	cli "modelhelper/cli/ports/app"
	"modelhelper/cli/ports/code"
	"modelhelper/cli/ports/command"
	"modelhelper/cli/ports/config"
	"modelhelper/cli/ports/converter"
	"modelhelper/cli/ports/language"
	"modelhelper/cli/ports/project"
	projectTemplate "modelhelper/cli/ports/project/template"
	"modelhelper/cli/ports/template"
	"os"
	"os/signal"
	"syscall"
)

// "fmt"
var version = "3.0.0-beta3"
var isBeta = true

func main() {

	// var (
	// 	// err  error
	// 	stop context.CancelFunc
	// )

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		os.Exit(3)
	// 	}
	// }()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// mainCtx := context.WithCancel()
	cmnd, ai := initializeApplication(ctx)

	if !ai.IsInitialized() {
		err := ai.Initialize()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	} else {
		cmnd.Execute(ctx)
	}
}

func initializeApplication(ctx context.Context) (modelhelper.CommandService, modelhelper.AppInitializer) {
	cfgs := config.NewConfigService()
	cfg, _ := cfgs.Load()

	prjs := project.NewProjectConfigService()
	infs := cli.NewCliInfo("modelhelper", version, isBeta)
	tpls := template.NewCodeTemplateService(cfg)
	lngs := language.NewLanguageDefinitionService(cfg)
	tcg := code.NewCodeGenerator(tpls, lngs)

	mha, _ := modelhelper.NewApplication(cfgs, prjs, infs)

	mha.Config = cfg
	mha.Code.TemplateService = tpls
	mha.Code.ModelConverter = converter.NewCodeModelConverter()
	mha.Code.Generator = code.NewCodeGeneratorService(cfg, mha.Project.Config, mha.Code.ModelConverter, tpls, tcg)

	mha.Project.TemplateService = projectTemplate.NewProjectTemplateService(cfg)
	mha.Project.Generator = project.NewProjectGeneratorService(cfg)
	mha.Project.ModelConverter = converter.NewProjectModelConverter()

	mha.LanguageService = lngs
	cmnd := command.NewCobraCli(mha)

	ai := cli.NewCliInitializer(ctx, infs, cfgs)
	return cmnd, ai
}

// func printTerminalSizes() {
// 	size, _ := ts.GetSize()
// 	fmt.Println(size.Col())  // Get Width
// 	fmt.Println(size.Row())  // Get Height
// 	fmt.Println(size.PosX()) // Get X position
// 	fmt.Println(size.PosY()) // Get Y position
// 	//
// }
