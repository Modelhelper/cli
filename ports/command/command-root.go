package command

import (
	"context"
	"fmt"
	"github.com/Modelhelper/cli/v3/modelhelper"
	"github.com/Modelhelper/cli/v3/ports/command/code"
	cfgCmd "github.com/Modelhelper/cli/v3/ports/command/config"
	"github.com/Modelhelper/cli/v3/ports/command/connection"
	"github.com/Modelhelper/cli/v3/ports/command/language"
	projectCmd "github.com/Modelhelper/cli/v3/ports/command/project"
	"github.com/Modelhelper/cli/v3/ports/command/serve"
	"github.com/Modelhelper/cli/v3/ports/command/source"
	"github.com/Modelhelper/cli/v3/ports/command/template"

	"github.com/spf13/cobra"
)

type cobraCommand struct {
	application *modelhelper.ModelhelperCli
	rootCmd     *cobra.Command
}

func NewCobraCli(application *modelhelper.ModelhelperCli) modelhelper.CommandService {

	return &cobraCommand{
		application: application,
		rootCmd:     createRootCommand(application),
	}
}

func createRootCommand(app *modelhelper.ModelhelperCli) *cobra.Command {
	return &cobra.Command{
		Use:   "mh",
		Short: "Shows information about the ModelHelper CLI",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(app.Info.Logo())
			fmt.Println(app.Info.About())

		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if f := cmd.Flags().Lookup("code-only"); f != nil && f.Value.String() == "true" {
				return
			}

			slogan := app.Info.Slogan()
			fmt.Printf("\n%s\n", slogan)
		},
	}
}

func (cc *cobraCommand) Execute(ctx context.Context) error {
	cc.buildCommandTree()

	if err := cc.rootCmd.ExecuteContext(ctx); err != nil {
		return err
	}

	return nil
}

func (cc *cobraCommand) buildCommandTree() {
	for _, cmd := range cc.subCommands() {
		cc.rootCmd.AddCommand(cmd)
	}
}

func (cc *cobraCommand) subCommands() []*cobra.Command {
	return []*cobra.Command{
		source.SourceCommand(cc.application),
		language.LanguageCommand(cc.application),
		projectCmd.ProjectCommand(cc.application),
		code.NewCodeRootCommand(cc.application),
		cfgCmd.NewConfigCommand(),
		template.NewTemplateCommand(cc.application),
		serve.NewServeCommand(),
		NewAboutCommand(cc.application.Info),
		NewVersionCommand(cc.application.Info),
		NewCompletionCommand(),
		connection.NewConnectionCommand(cc.application),
	}

}
