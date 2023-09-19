package language

import (
	"modelhelper/cli/modelhelper"

	"github.com/spf13/cobra"
)

func LanguageCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	subCommands := []*cobra.Command{
		NewListLanguagesCommand(app),
	}

	rootCmd := &cobra.Command{
		Use:   "language",
		Short: "Manage modelhelper languages",
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}
