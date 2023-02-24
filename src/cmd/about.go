package cmd

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/ports/project"

	"github.com/spf13/cobra"
)

func NewAboutCommand(info modelhelper.AppInfoService) *cobra.Command {

	// aboutCmd represents the about command
	return &cobra.Command{
		Use:   "about",
		Short: "Show information about the modelhelper CLI",
		Run: func(cmd *cobra.Command, args []string) {
			printLogoInfo(info)
		},
	}
}

func printLogoInfo(info modelhelper.AppInfoService) {
	fmt.Print(info.Logo())

	dir := project.DefaultLocation()

	if project.Exists(dir) {
		// printProjectInfo(project.DefaultLocation(), true)

	} else {
		fmt.Println(info.About())
	}

}
