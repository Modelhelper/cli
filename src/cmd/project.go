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
	Short: "Working with the projects",

	Run: func(cmd *cobra.Command, args []string) {
		open, _ := cmd.Flags().GetBool("open")

		if project.Exists(project.DefaultLocation()) {

			if open {
				openProjectInEditor()
			} else {
				printProjectInfo(project.DefaultLocation(), true)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.Flags().Bool("open", false, "Opens the project file in default editor")

}

func printProjectInfo(projectFile string, renderTables bool) {
	p, err := project.Load(projectFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg := config.Load()
	cons := source.JoinConnections("smart", cfg, p)

	ui.PrintConsoleTitle("Project information")

	fmt.Printf("\n%-20s%8s", "Name", p.Name)
	fmt.Printf("\n%-20s%8s", "Version", p.Version)
	fmt.Printf("\n%-20s%8s", "Owner", p.OwnerName)
	fmt.Printf("\n%-20s%8s", "Primary language", p.Language)
	fmt.Printf("\n\n")

	// fmt.Println("Defaults:")
	// fmt.Printf("%-20s%8s", "Connection", p.DefaultSource)
	// fmt.Printf("\n%-20s%8s", "Key", p.DefaultKey)

	if len(cons) > 0 {
		fmt.Printf("\n\n")
		fmt.Println("Available Connections:")
		fmt.Printf("\n")

		cr := connectionRenderer{cons, p.DefaultSource}
		ui.RenderTable(&cr, &cr)
		fmt.Printf("\n")
	}
	if renderTables {

		if len(p.Code.Keys) > 0 {

			kr := keyRenderer{keys: p.Code.Keys}
			ui.RenderTable(&kr, &kr)
		}

		if len(p.Code.Inject) > 0 {

			ir := injectRenderer{p.Code.Inject}
			ui.RenderTable(&ir, &ir)
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

type connectionRenderer struct {
	rows   map[string]source.Connection
	defcon string
}

func (l *connectionRenderer) BuildHeader() []string {
	return []string{
		"Key",
		"Default",
		"Type",
		"Description",
	}
}

func (d *connectionRenderer) ToRows() [][]string {
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
