/*
Copyright © 2021 Hans-Petter Eitvet

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"modelhelper/cli/config"
	"modelhelper/cli/project"
	"modelhelper/cli/ui"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// projectInitCmd represents the projectInit command
var projectInitCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new project in the current working directory",

	Run: func(cmd *cobra.Command, args []string) {
		p := *&project.Project{}

		init := true

		if project.Exists(project.DefaultLocation()) {
			color.Red.Println("NB!!")
			color.Red.Println("A project already exists in this location")
			init = ui.PromptForYesNo("Overwrite current project file? [y/N]", "n")
		}

		if init {
			p.Version = "3.0"
			p.Name = ui.PromptForString("Enter the name of the project")
			p.DefaultSource = promptForConnectionKey()
			p.Language = ui.PromptForLanguage("Select the primary code language")
			p.Options = make(map[string]string)
			// p.OwnerName = promptForString("Enter the owner (company name) for this project")

			if ui.PromptForYesNo("Clone connections from config? [Y/n]", "y") {
				cfg := config.Load()
				p.Connections = cfg.Connections
				// clone
			}

			err := p.Save()

			if err != nil {
				fmt.Println(err)
			}

			open, _ := cmd.Flags().GetBool("open")

			if project.Exists(project.DefaultLocation()) {

				if open {
					openProjectInEditor()
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(projectInitCmd)
	projectCmd.AddCommand(projectInitCmd)
	projectInitCmd.Flags().Bool("open", false, "Opens the project file in default editor")

}
