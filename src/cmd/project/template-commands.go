package project

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"modelhelper/cli/ui"
	"strings"

	"github.com/spf13/cobra"
)

func NewTemplatesCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "template",
		Short: "Creates a new project in the current working directory",

		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("mh project template")
		},
	}

	cmd.AddCommand(templateListCommand(app))
	cmd.AddCommand(templateDetailCommand(app))

	return cmd
}

func templateListCommand(app *modelhelper.ModelhelperCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Creates a new project in the current working directory",

		Run: func(cmd *cobra.Command, args []string) {

			ui.PrintConsoleTitle("Project Templates")
			ls := app.Project.TemplateService.List(nil)

			presenter := &templatePrinter{ls}

			ui.RenderTable(presenter)

		},
	}

	return cmd
}

func templateDetailCommand(app *modelhelper.ModelhelperCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detail [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Get details about a named template",

		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("mh project template detail")
			tpl := app.Project.TemplateService.Load(args[0])

			fmt.Printf("name: %s, lang: %s", tpl.Name, tpl.Language)
		},
	}

	return cmd
}

type templatePrinter struct {
	templates map[string]models.ProjectTemplate
}

func (t *templatePrinter) Rows() [][]string {
	var rows [][]string

	for name, t := range t.templates {
		tags := strings.Join(t.Tags, ", ")
		row := []string{
			name,
			t.Language,
			// t.Model,
			// t.Key,
			tags,
			// t.Short,
			t.Description,
		}

		rows = append(rows, row)
	}

	return rows
}

func (t *templatePrinter) Header() []string {

	row := []string{
		"Name",
		"Language",
		// "Type",
		// "Model",
		// "Key",
		"Tags",
		"Description",
	}

	return row
}
