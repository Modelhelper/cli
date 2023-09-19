package modelhelper

import "modelhelper/cli/modelhelper/models"

type ModelhelperCli struct {
	Config        *models.Config
	ConfigService ConfigService

	Project struct {
		Gen             TemplateGenerator[*models.ProjectTemplate]
		Exists          bool
		Config          *models.ProjectConfig
		ConfigService   ProjectConfigService
		TemplateService ProjectTemplateService
		Generator       ProjectGenerator
		ModelConverter  ProjectModelConverter
	}

	Code struct {
		TemplateService CodeTemplateService
		Generator       CodeGeneratorService
		ModelConverter  CodeModelConverter
		CommitHistory   CommitHistoryService
	}

	Exporters struct {
		ScreenExporter    Exporter
		ClipboardExporter Exporter
		FileExporter      Exporter
	}

	ConnectionService ConnectionService
	SourceFactory     SourceFactoryService
	LanguageService   LanguageDefinitionService
	Version           string
	IsBeta            bool
	Info              AppInfoService
}

func NewApplication(cfgService ConfigService, projectService ProjectConfigService, infoService AppInfoService) (*ModelhelperCli, error) {
	app := &ModelhelperCli{
		ConfigService: cfgService,
		Info:          infoService,
	}

	app.Project.ConfigService = projectService

	p, err := projectService.Load()
	if err != nil {
		app.Project.Exists = false
	}

	app.Project.Exists = p != nil
	app.Project.Config = p

	return app, nil
}
