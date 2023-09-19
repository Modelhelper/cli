package command

import (
	"fmt"
	"modelhelper/cli/modelhelper"

	"github.com/spf13/cobra"
)

func NewVersionCommand(info modelhelper.AppInfoService) *cobra.Command {

	// versionCmd represents the version command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the CLI version information",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(info.Version())
		},
	}

	return versionCmd
}
