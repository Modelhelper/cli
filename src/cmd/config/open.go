package config

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

func NewOpenConfigCommand() *cobra.Command {

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
	fmt.Println("config open")
}

func openPathInEditor(editor string, loc string) {
	exe := exec.Command(editor, loc)
	if exe.Run() != nil {
		//vim didn't exit with status code 0
	}
}
