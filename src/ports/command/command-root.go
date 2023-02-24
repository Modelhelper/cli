package command

import (
	"context"
	"fmt"
	"modelhelper/cli/cmd"
	"modelhelper/cli/cmd/code"
	cfgCmd "modelhelper/cli/cmd/config"
	"modelhelper/cli/cmd/language"
	projectCmd "modelhelper/cli/cmd/project"
	"modelhelper/cli/cmd/serve"
	"modelhelper/cli/cmd/source"
	"modelhelper/cli/cmd/template"
	"modelhelper/cli/modelhelper"

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
		source.SourceCommand(),
		language.LanguageCommand(cc.application),
		projectCmd.ProjectCommand(cc.application),
		code.NewCodeRootCommand(cc.application),
		cfgCmd.NewConfigCommand(),
		template.NewTemplateCommand(cc.application),
		serve.NewServeCommand(),
		cmd.NewAboutCommand(cc.application.Info),
		cmd.NewVersionCommand(cc.application.Info),
		cmd.NewCompletionCommand(),
	}

}
