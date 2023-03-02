package template

import (
	"modelhelper/cli/modelhelper"

	"github.com/spf13/cobra"
)

func NewTemplateCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	subCommands := []*cobra.Command{
		ListCommand(app),
		OpenCommand(app),
		CreateCommand(app),
	}

	rootCmd := &cobra.Command{
		Use:     "template",
		Aliases: []string{"t"},
		Short:   "Manage modelhelper templates",
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}
