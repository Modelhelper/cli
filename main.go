package main

import (
	"context"
	"fmt"
	"log/slog"
	"modelhelper/cli/modelhelper"
	cli "modelhelper/cli/ports/app"
	"modelhelper/cli/ports/code"
	"modelhelper/cli/ports/command"
	"modelhelper/cli/ports/config"
	"modelhelper/cli/ports/connection"
	"modelhelper/cli/ports/converter"
	"modelhelper/cli/ports/language"
	"modelhelper/cli/ports/project"
	projectTemplate "modelhelper/cli/ports/project/template"
	"modelhelper/cli/ports/source"
	"modelhelper/cli/ports/template"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	_ "embed"
)

//go:embed version.txt
var version string

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
	defer handlePanic()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfgs := config.NewConfigService()
	// mainCtx := context.WithCancel()
	if !cfgs.ConfigExists() {
		cfg, err := cfgs.Create()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		// fmt.Println("Config created", cfg)
		err = cfgs.SaveConfig(cfg)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	cmnd, ai := initializeApplication(ctx, cfgs)

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

func handlePanic() {
	if err := recover(); err != nil {
		slog.Error("context", "main", "message", "Caught panic", "error", err)
		slog.Error(string(debug.Stack()))
	}
}

func initializeApplication(ctx context.Context, cfgs modelhelper.ConfigService) (modelhelper.CommandService, modelhelper.AppInitializer) {
	cfg, _ := cfgs.Load()

	prjs := project.NewProjectConfigService()
	infs := cli.NewCliInfo("modelhelper", version, isBeta)

	tpls := template.NewCodeTemplateService(cfg, prjs.TemplatePath())

	lngs := language.NewLanguageDefinitionService(cfg)
	tcg := code.NewCodeGenerator(tpls, lngs)
	cons := connection.NewConnectionService(cfg)
	srcs := source.NewSourceFactoryService(cfg, cons)
	mha, _ := modelhelper.NewApplication(cfgs, prjs, infs)

	mha.Config = cfg

	cmtHist := code.NewCommitHistoryService()
	mha.Code.TemplateService = tpls
	mha.Code.ModelConverter = converter.NewCodeModelConverter()
	mha.Code.CommitHistory = cmtHist
	mha.Code.Generator = code.NewCodeGeneratorService(cfg, mha.Project.Config, mha.Code.ModelConverter, tpls, tcg, cons, srcs, cmtHist)

	mha.Project.TemplateService = projectTemplate.NewProjectTemplateService(cfg)
	mha.Project.Generator = project.NewProjectGeneratorService(cfg)
	mha.Project.ModelConverter = converter.NewProjectModelConverter()
	mha.ConnectionService = cons
	mha.SourceFactory = srcs
	mha.LanguageService = lngs
	cmnd := command.NewCobraCli(mha)

	ai := cli.NewCliInitializer(ctx, infs, cfgs, cons)
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

/*

	//%%MH_SNIPPET_QUERY_DEF%%
GetCalibration: calibration.NewGetCalibrationHandler(sqlConn, logger),
GetCalibration: calibration.NewGetCalibrationHandler(sqlConn, logger),



*/
