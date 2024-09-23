package project

import (
	"encoding/json"
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/ports/converter"

	"github.com/spf13/cobra"
)

func ProjectCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	subCommands := []*cobra.Command{
		NewGenerateProjectCommand(app),
		NewOpenProjectCommand(app),
		NewProjectInitCommand(app),
		NewTemplatesCommand(app),
	}

	rootCmd := &cobra.Command{
		Use:     "project",
		Aliases: []string{"p"},
		Short:   "Work with projects",
		Run: func(cmd *cobra.Command, args []string) {
			model, _ := cmd.Flags().GetBool("model")
			if model {
				writeProjectModel(app)
			} else {
				writeProjectInfo(app)
			}
		},
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	rootCmd.Flags().Bool("model", false, "Show project as a basic template model")

	// --model flag to view the project model || basic model
	return rootCmd
}

func writeProjectInfo(app *modelhelper.ModelhelperCli) {
	if !app.Project.Exists {
		fmt.Printf("No project exists here \n")
		return
	}

	fmt.Printf("Name: %s\n", app.Project.Config.Name)
	fmt.Printf("Description: %s\n", app.Project.Config.Description)
	fmt.Printf("Language: %s\n", app.Project.Config.Language)
	fmt.Printf("Root namespace: %s\n", app.Project.Config.RootNamespace)
	fmt.Printf("Owner: %s\n", app.Project.Config.OwnerName)

	fmt.Printf("\nFeatures\n")

	useAuth := app.Project.Config.Features != nil && app.Project.Config.Features.Auth != nil && app.Project.Config.Features.Auth.Use
	useLogger := app.Project.Config.Features != nil && app.Project.Config.Features.Logger != nil && app.Project.Config.Features.Logger.Use
	useTracing := app.Project.Config.Features != nil && app.Project.Config.Features.Tracing != nil && app.Project.Config.Features.Tracing.Use
	useSwagger := app.Project.Config.Features != nil && app.Project.Config.Features.Swagger != nil && app.Project.Config.Features.Swagger.Use
	useMetrics := app.Project.Config.Features != nil && app.Project.Config.Features.Metrics != nil && app.Project.Config.Features.Metrics.Use
	useHealth := app.Project.Config.Features != nil && app.Project.Config.Features.Health != nil && app.Project.Config.Features.Health.Use
	useApi := app.Project.Config.Features != nil && app.Project.Config.Features.Api != nil && app.Project.Config.Features.Api.Use
	useDb := app.Project.Config.Features != nil && app.Project.Config.Features.Db != nil && app.Project.Config.Features.Db.Use

	fmt.Printf("\tAuth:\t\t%v", useAuth)
	fmt.Printf("\n\tLogger:\t\t%v", useLogger)
	fmt.Printf("\n\tTracing:\t%v\n", useTracing)
	fmt.Printf("\tSwagger:\t%v\n", useSwagger)
	fmt.Printf("\tMetrics:\t%v\n", useMetrics)
	fmt.Printf("\tHealth:\t\t%v\n", useHealth)
	fmt.Printf("\tApi:\t\t%v\n", useApi)
	fmt.Printf("\tDb:\t\t%v\n", useDb)

	for k, v := range app.Project.Config.CustomFeatures {
		fmt.Printf("%s: %s\n", k, v.Namespace)
	}

	fmt.Printf("\nSetup\n")
	for k, v := range app.Project.Config.Setup {
		fmt.Printf("\t%s: {namespace: %s, prefix: %s, postfix: %s}\n", k, v.Namespace, v.Prefix, v.Postfix)
	}

	fmt.Printf("\nCode Export Locations\n")
	if app.Project.Config.Locations == nil && len(app.Project.Config.Locations) == 0 {
		fmt.Printf("\tNo locations set\n")
	} else {
		for k, v := range app.Project.Config.Locations {
			fmt.Printf("\t%s: %s\n", k, v)
		}
	}

	if app.Project.Config.Inject != nil && len(app.Project.Config.Inject) > 0 {

		fmt.Printf("\nInject\n")
		for k, v := range app.Project.Config.Inject {
			fmt.Printf("\t%s: {property: %s, name: %s, method: %s}\n", k, v.PropertyName, v.Name, v.Method)
		}
	}

	if app.Project.Config.Options != nil && len(app.Project.Config.Options) > 0 {

		fmt.Printf("\nOptions\n")
		for k, v := range app.Project.Config.Options {
			fmt.Printf("\t%s: %s\n", k, v)
		}
	}
}

func writeProjectModel(app *modelhelper.ModelhelperCli) {
	if !app.Project.Exists {
		fmt.Printf("No project exists here \n")
		return
	}

	conv := converter.NewCodeModelConverter()

	id := "model"
	proj := app.Project.Config
	lang := proj.Language
	model := conv.ToBasicModel(id, lang, proj)

	jsonB, _ := json.MarshalIndent(model, "", "  ")

	fmt.Printf("Project as model\n\n %v\n", string(jsonB))
}
