package cmd

import (
	"fmt"
	"modelhelper/cli/app"
	_ "modelhelper/cli/types"

	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {

	// versionCmd represents the version command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the CLI version information",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(app.Version())
		},
	}

	return versionCmd
}
