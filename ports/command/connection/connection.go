package connection

import (
	"modelhelper/cli/modelhelper"

	"github.com/spf13/cobra"
)

func NewConnectionCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	subCommands := []*cobra.Command{
		NewListConnectionsCommand(app),
		// NewCreateConnectionCommand(app.ConnectionService),
		// NewDeleteConnectionCommand(app.ConnectionService),
	}

	rootCmd := &cobra.Command{
		Use:   "connection",
		Short: "Manage modelhelper connections",
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}
