package connection

import (
	"fmt"
	"modelhelper/cli/modelhelper"

	"github.com/spf13/cobra"
)

func NewCreateConnectionCommand(cs modelhelper.ConnectionService) *cobra.Command {
	cmd := &cobra.Command{
		Example: "mh connection create mssql 'name'",
		Use:     "create [mssql|file] {name}",
		Aliases: []string{"c"},
		// Args:    cobra.MaximumNArgs(1),
		Short: "Create a connections, use sub commands for each specific type of connection",
		Long: `
This is the root command for creating a connection. You must combine this
command with a sub command of either 'mssql' or 'file', followed by the name of the connection

Use the various options for each connection type to give details.`,
		// Run:
	}

	cmd.AddCommand(NewCreateMsSQLConnectionCommand(cs))
	cmd.AddCommand(NewCreateFileConnectionCommand(cs))
	return cmd
}

func NewCreateMsSQLConnectionCommand(cs modelhelper.ConnectionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mssql",
		Args:  cobra.MinimumNArgs(1),
		Short: "Creates a MS SQL connection",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Create a file connection")
		},
	}

	return cmd
}

func NewCreateFileConnectionCommand(cs modelhelper.ConnectionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file",
		Args:  cobra.MinimumNArgs(1),
		Short: "Creates a file connection",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Create a file connection")
		},
	}

	return cmd
}
