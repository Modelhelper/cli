package cmd

import (
	"fmt"
	"modelhelper/cli/app"
	cfCmd "modelhelper/cli/cmd/config"
	"modelhelper/cli/cmd/language"
	prCmd "modelhelper/cli/cmd/project"
	"modelhelper/cli/cmd/serve"
	"modelhelper/cli/cmd/source"
	tplCmd "modelhelper/cli/cmd/template"
	"os"

	"github.com/spf13/cobra"
)

var modelHelperApp *app.Application

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mh",
	Short: "Shows information about the ModelHelper CLI",

	Run: func(cmd *cobra.Command, args []string) {

		printLogoInfo()

	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		a := app.NewModelhelperCli()
		slog := a.Slogan()
		fmt.Println(slog)
	},
}

func SetApplication(app *app.Application) {
	modelHelperApp = app
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {

		fmt.Println(err)
		os.Exit(1)
	}

}

func subCommands() []*cobra.Command {
	return []*cobra.Command{
		source.SourceCommand(),
		language.LanguageCommand(),
		prCmd.ProjectCommand(),
		cfCmd.NewConfigCommand(),
		tplCmd.NewTemplateCommand(),
		serve.NewServeCommand(),
		NewAboutCommand(),
		NewVersionCommand(),
		NewCompletionCommand(),
	}

}

func init() {

	for _, cmd := range subCommands() {
		rootCmd.AddCommand(cmd)
	}
	// cobra.OnInitialize(initConfig)
	//rootCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "Sets the source")

}
