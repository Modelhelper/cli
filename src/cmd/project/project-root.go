package project

import (
	"fmt"
	"modelhelper/cli/modelhelper"

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
		Short:   "Generates code based on language, template and source",
		Run:     rootCommandHandler,
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}

func rootCommandHandler(cmd *cobra.Command, args []string) {
	fmt.Println("project info")
}
