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
		Short:   "Work with projects",
		Run: func(cmd *cobra.Command, args []string) {
			writeProjectInfo(app)
		},
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}

func writeProjectInfo(app *modelhelper.ModelhelperCli) {
	if !app.Project.Exists {
		fmt.Printf("No project exists here \n")
		return
	}

	fmt.Printf(app.Project.Config.Name)
}
