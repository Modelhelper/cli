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

	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// projectInitCmd represents the projectInit command
var projectInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new project in the current working directory",

	Run: func(cmd *cobra.Command, args []string) {
		p := *&project.Project{}

		init := true

		if project.Exists(project.DefaultLocation()) {
			color.Red.Println("NB!!")
			color.Red.Println("A project already exists in this location")
			init = promptForYesNo("Overwrite current project file? ")
		}

		if init {
			p.Version = "3.0"
			p.Name = promptForString("Enter the name of the project")
			p.DefaultSource = promptForConnectionKey()
			p.Language = promptForLanguage()
			p.Options = make(map[string]string)
			// p.OwnerName = promptForString("Enter the owner (company name) for this project")

			if promptForYesNo("Clone connections from config? ") {
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

func promptForString(question string) string {

	prompt := promptui.Prompt{
		Label: question,
	}

	result, err := prompt.Run()

	if err != nil {
		return ""
	}

	return result
}

func promptForYesNo(question string) bool {
	answers := make(map[string]bool)
	answers["Yes"] = true
	answers["No"] = false

	items := []string{}

	for k, _ := range answers {
		items = append(items, k)
	}

	index := -1
	var key string
	var err error

	for index < 0 {
		prompt := promptui.Select{
			Label: question,
			Items: items,
		}

		index, key, err = prompt.Run()

		if index == -1 {
			items = append(items, key)
		}
	}

	if err != nil {
		// fmt.Printf("Prompt failed %v\n", err)
		return false
	}

	selected, found := answers[key]
	if found {
		return selected
	}

	return false
}
func promptForLanguage() string {
	lang := make(map[string]string)
	lang["C#"] = "cs"
	lang["Go"] = "go"
	lang["TypeScript"] = "ts"
	lang["JavaScript"] = "js"
	lang["Python"] = "py"
	lang["Java"] = "java"

	items := []string{}

	for k, _ := range lang {
		items = append(items, k)
	}

	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    "Select the primary code language",
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

	selected, found := lang[result]
	if found {
		return selected
	}

	return result
}
