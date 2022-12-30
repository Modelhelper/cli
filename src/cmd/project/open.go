package project

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewOpenProjectCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "open",
		Short: "Opens the project config file in an editor",
		Long:  "",
		Run:   openCommand,
	}

	cmd.Flags().String("editor", "", "The editor to use when opening the file")

	return cmd
}

func openCommand(cmd *cobra.Command, args []string) {
	fmt.Println("project open")
}
