package connection

import (
	"fmt"
	"modelhelper/cli/modelhelper"

	"github.com/spf13/cobra"
)

func NewDeleteConnectionCommand(cs modelhelper.ConnectionService) *cobra.Command {
	cmd := &cobra.Command{
		Example: "mh connection delete 'name'",
		Use:     "delete [name]",
		Aliases: []string{"d"},
		// Args:    cobra.MaximumNArgs(1),
		Short: "Deletes a named connection",
		Run: func(cmd *cobra.Command, args []string) {
			con := "named"
			fmt.Printf("Deletes a %s connection", con)
		},
	}
	return cmd
}
