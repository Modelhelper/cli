package cmd

import (
	"fmt"
	"modelhelper/cli/app"
	"modelhelper/cli/project"

	"github.com/gookit/color"

	"github.com/spf13/cobra"
)

func NewAboutCommand() *cobra.Command {

	// aboutCmd represents the about command
	return &cobra.Command{
		Use:   "about",
		Short: "Show information about the modelhelper CLI",
		Run: func(cmd *cobra.Command, args []string) {
			printLogoInfo()
		},
	}
}

// func init() {
// 	rootCmd.AddCommand(aboutCmd)
// }

func printLogoInfo() {
	a := app.NewModelhelperCli()
	color.Green.Print(a.Logo())

	dir := project.DefaultLocation()

	if project.Exists(dir) {
		// printProjectInfo(project.DefaultLocation(), true)

	} else {
		fmt.Println(a.About())
	}

}
