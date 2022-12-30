package template

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewOpenTemplateCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "open",
		Short: "Opens the template file in an editor",
		Long:  "",
		Run:   openCommand,
	}

	cmd.Flags().String("editor", "", "The editor to use when opening the file")

	return cmd
}

func openCommand(cmd *cobra.Command, args []string) {
	fmt.Println("project open")
}
