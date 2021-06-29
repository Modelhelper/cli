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
	"fmt"
	"modelhelper/cli/config"
	"modelhelper/cli/project"
	"modelhelper/cli/source"
	"modelhelper/cli/ui"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage the current project in the working directory",

	Run: func(cmd *cobra.Command, args []string) {
		open, _ := cmd.Flags().GetBool("open")
		path, _ := cmd.Flags().GetString("path")
		// projects := project.LoadProjects(project.FindReleatedProjects()...)

		if project.Exists(path) {

			if open {
				openProjectInEditor()
			} else {
				printProjectInfo(path, true)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.Flags().Bool("open", false, "Opens the project file in default editor")
	projectCmd.Flags().String("path", project.DefaultLocation(), "Opens the project file in default editor")

}

func printProjectInfo(projectFile string, renderTables bool) {

	paths := project.FindReleatedProjects(projectFile)
	ps := project.LoadProjects(paths...)

	p := project.JoinProject("smart", ps...)
	// p, err := project.Load(projectFile)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	cfg := config.Load()
	cons := source.JoinConnections("smart", cfg, &p)

	ui.PrintConsoleTitle("Project information")

	fmt.Printf("\n%-20s%20s", "Name", p.Name)
	// fmt.Printf("\n%-20s%20s", "Version", p.Version)
	fmt.Printf("\n%-20s%20s", "Owner", p.OwnerName)
	fmt.Printf("\n%-20s%20s", "Primary language", p.Language)
	fmt.Printf("\n\n")
	fmt.Println(p.Description)

	fmt.Println()
	fmt.Println("OPTIONS....")
	for k, v := range p.Options {
		fmt.Printf("\n%-20s%20s", k, v)
	}
	// fmt.Println("Defaults:")
	// fmt.Printf("%-20s%8s", "Connection", p.DefaultSource)
	// fmt.Printf("\n%-20s%8s", "Key", p.DefaultKey)

	if len(cons) > 0 {
		fmt.Printf("\n\n")
		fmt.Println("Available Connections:")
		fmt.Printf("\n")

		cr := connectionRenderer{cons, p.DefaultSource}
		ui.RenderTable(&cr)
		fmt.Printf("\n")
	}
	if renderTables {
		for _, langVal := range p.Code {

			showTemplateKey := false
			if len(langVal.Keys) > 0 {
				ui.PrintConsoleTitle("Keys")
				kr := keyRenderer{keys: langVal.Keys}
				ui.RenderTable(&kr)
			} else {
				fmt.Printf(`No keys is defined for this project
Using keys will enable the templates to render correct namespace, package etc

use the command 'mh project key <name> --namespace 'namespace'
`)

				showTemplateKey = true
			}

			if len(langVal.Inject) > 0 {
				ui.PrintConsoleTitle("Inject")
				ir := injectRenderer{langVal.Inject}
				ui.RenderTable(&ir)
			} else {
				showTemplateKey = true
			}
			if len(langVal.Locations) > 0 {
				ui.PrintConsoleTitle("Locations")

				lr := locationRenderer{langVal.Locations}
				ui.RenderTable(&lr)

			} else {
				fmt.Printf(`Locations is not defined for this project
Connecting keys to a location will enable the templates export generated code to a path relative to this project

use the command 'mh project location <name> <path>
`)
				showTemplateKey = true
			}

			if showTemplateKey {
				fmt.Println()
				fmt.Println("Where to find keys to use in project")
				fmt.Println("use 'mh template' to see which keys that each template implement")
			}
		}
	}
}

func openProjectInEditor() {
	cfg := config.Load()
	editor := cfg.DefaultEditor

	if len(cfg.DefaultEditor) == 0 {
		editor = promptForEditorKey("Please select a editor to open the project file")
	}

	openPathInEditor(editor, project.DefaultLocation())
}

type locationRenderer struct {
	rows map[string]string
}

func (l *locationRenderer) Header() []string {
	return []string{
		"Key",
		"Path",
	}
}

func (d *locationRenderer) Rows() [][]string {
	var rows [][]string
	// p := message.NewPrinter(language.English)

	for key, val := range d.rows {

		r := []string{
			key,
			val,
		}

		rows = append(rows, r)
	}

	return rows
}

type connectionRenderer struct {
	rows   map[string]source.Connection
	defcon string
}

func (l *connectionRenderer) Header() []string {
	return []string{
		"Key",
		"Default",
		"Type",
		"Description",
	}
}

func (d *connectionRenderer) Rows() [][]string {
	var rows [][]string
	// p := message.NewPrinter(language.English)

	for key, val := range d.rows {
		def := ""

		if d.defcon == key {
			def = "Yes"
		}
		r := []string{
			key,
			def,
			val.Type,
			val.Description,
		}

		rows = append(rows, r)
	}

	return rows
}
