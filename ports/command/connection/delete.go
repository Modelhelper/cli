package connection

import (
	"fmt"
	"modelhelper/cli/modelhelper"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func NewDeleteConnectionCommand(cs modelhelper.ConnectionService) *cobra.Command {

	var (
		confirmed bool
	)
	cmd := &cobra.Command{
		Example: "mh connection delete 'name'",
		Use:     "delete [name]",
		Aliases: []string{"del", "rm"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Deletes a named connection",
		Run: func(cmd *cobra.Command, args []string) {
			var connectionName string
			var err error

			if len(args) == 0 {
				connectionName, err = getConnectionNameFromSelectionList(cs)
				if err != nil {
					fmt.Println("Failed to get connection: ", err)
					return
				}
			}

			if len(args) > 0 {
				if !confirmed {
					fmt.Println("Please confirm the deletion")
					return
				}

				connectionName = args[0]

			}

			if !confirmed {

				confirm := huh.NewConfirm().
					// Title("Are you sure? ").
					TitleFunc(func() string { return "Are you sure to delete " + connectionName }, &connectionName).
					// Description("Please confirm. ").
					Affirmative("Yes").
					Negative("Cancel").
					Value(&confirmed)

				huh.NewForm(huh.NewGroup(confirm)).Run()
			}

			if confirmed {

				err := cs.Delete(connectionName)
				if err != nil {
					fmt.Println("Failed to delete connection: ", err)
					return
				}

				fmt.Print("Connection deleted")
			}
		},
	}

	cmd.Flags().BoolVar(&confirmed, "confirm", false, "Confirms the deletion")
	// cmd.Flags().StringV("repo", "r", "", "Where the repo is located")

	return cmd
}
