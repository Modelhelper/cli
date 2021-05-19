/*
Copyright © 2020 Hans-Petter Eitvet

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
	"modelhelper/cli/config"
	"modelhelper/cli/project"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Working with the projects",

	Run: func(cmd *cobra.Command, args []string) {
		open, _ := cmd.Flags().GetBool("open")

		if project.Exists(project.DefaultLocation()) {

			if open {
				openProjectInEditor()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.Flags().Bool("open", false, "Opens the project file in default editor")

}

func openProjectInEditor() {
	cfg := config.Load()
	editor := cfg.DefaultEditor

	if len(cfg.DefaultEditor) == 0 {
		editor = promptForEditorKey("Please select a editor to open the project file")
	}

	openPathInEditor(editor, project.DefaultLocation())
}
