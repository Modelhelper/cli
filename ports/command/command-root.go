package command

import (
	"context"
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/ports/command/code"
	cfgCmd "modelhelper/cli/ports/command/config"
	"modelhelper/cli/ports/command/connection"
	"modelhelper/cli/ports/command/language"
	projectCmd "modelhelper/cli/ports/command/project"
	"modelhelper/cli/ports/command/serve"
	"modelhelper/cli/ports/command/source"
	"modelhelper/cli/ports/command/template"

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
			slog := app.Info.Slogan()
			fmt.Print(slog)
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
