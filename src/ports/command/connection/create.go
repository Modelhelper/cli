package connection

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"

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
			cons, _ := cs.Connections()
			name := args[0]

			_, exists := cons[name]

			if exists {
				panic("A connection with the name already exists")
			}

			con := &models.MsSqlConnection{}

			cs, _ := cmd.Flags().GetString("connection-string")

			if len(cs) > 0 {
				con.ConnectionString = cs

			} else {

				s, _ := cmd.Flags().GetString("server")
				d, _ := cmd.Flags().GetString("database")
				u, _ := cmd.Flags().GetString("user")
				p, _ := cmd.Flags().GetString("password")

				if len(s) > 0 && len(d) > 0 && len(u) > 0 && len(p) > 0 {
					con.Database = d
					con.Server = s
					con.ConnectionString = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", u, p, s, d)
				}
			}
		},
	}

	cmd.Flags().StringP("connection-string", "c", "", "The full connection string to use")
	cmd.Flags().StringP("server", "s", "", "The server")
	cmd.Flags().StringP("database", "d", "", "The database")
	cmd.Flags().StringP("username", "u", "", "The user")
	cmd.Flags().StringP("password", "p", "", "The password")
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
