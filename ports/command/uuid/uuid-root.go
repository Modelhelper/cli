package uuid

import (
	"modelhelper/cli/modelhelper"

	"github.com/spf13/cobra"
)

func NewUUIDRootCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	subCommands := []*cobra.Command{
		NewGenerateUUIDCommand(app),
	}

	rootCmd := &cobra.Command{
		Use:     "uuid",
		Aliases: []string{"c"},
		Short:   "Generates UUIDs",
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}
