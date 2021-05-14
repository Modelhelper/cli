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
	"modelhelper/cli/config"

	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// setDefaultEditorCmd represents the setDefaultEditor command
var setDefaultEditorCmd = &cobra.Command{
	Use:  "editor",
	Args: cobra.MaximumNArgs(1),

	Short: "Sets the default editor",
	Long: `Sets the default editor
This is for opening config, project, templates and so on in the correct editor

The value should be a valid handle for opening the editor from a terminal 
e.g code for VSCode...

`,
	Run: func(cmd *cobra.Command, args []string) {

		editor := ""
		if len(args) == 0 {
			editor = promptForEditorKey()
		} else {
			editor = args[0]
		}

		err := config.SetDefaultEditor(editor)

		if err != nil {
			color.Red.Println(err)
			return
		}

		color.Green.Printf("\nSuccessfully updated the default editor to %s\n", editor)

	},
}

func init() {
	setCmd.AddCommand(setDefaultEditorCmd)

}

func promptForEditorKey() string {
	items := []string{"Vim", "Emacs", "Sublime", "VSCode", "Atom"}
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    "What's your editor",
			Items:    items,
			AddLabel: "Other",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		// fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return result
}
