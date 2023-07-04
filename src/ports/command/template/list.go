package template

import (
	"fmt"
	"modelhelper/cli/modelhelper"

	"modelhelper/cli/modelhelper/models"
	"modelhelper/cli/ui"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func ListCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "List all templates or view content of a single template",
		Long: `With this command you can list all available templates, snippet and blocks
		
		with the use of --by option you can group the various templates either by type, language
		group or tag.
		
A template can exist in one or more groups making it possible to generate different
outcome based on what you need

Filtering the templates:
Filter the template by using on or more of the following options
--lang <langcodes> (e.g --lang cs), filters by language
--type <type> (e.g --type block), filters by type
--model <model> (e.g --model entity), filters by model
--group <groupname> (e.g --group cs-dpr-full), filters by group


-- hva med liste på gruppenavn
-- hva med liste på tags

`,
		Run: listTemplateCommandHandler(app),
	}

	cmd.Flags().String("db", "pg", "Use a specific database type to get templates for this database [mssql, pg, postgres]")
	cmd.Flags().String("by", "", "Groups the templates by type, group, language, model or tag")
	cmd.Flags().StringArray("type", []string{}, "Filter the templates by the name of the type")
	cmd.Flags().StringArray("lang", []string{}, "Filter the templates by language")
	cmd.Flags().StringArray("model", []string{}, "Filter the templates by model")
	cmd.Flags().StringArray("key", []string{}, "Filter the templates by key")
	cmd.Flags().StringArray("group", []string{}, "Filter the templates by group")

	// templateCmd.Flags().Bool("open", false, "Opens the template file in default editor or a selection of editors")
	// 	cmd.Flags().String("editor", "", "The editor to use when opening the file")

	return cmd
}

func listTemplateCommandHandler(app *modelhelper.ModelhelperCli) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		db, _ := cmd.Flags().GetString("db")
		group, _ := cmd.Flags().GetString("by")
		typeFiler, _ := cmd.Flags().GetStringArray("type")
		langFilter, _ := cmd.Flags().GetStringArray("lang")
		modelFilter, _ := cmd.Flags().GetStringArray("model")
		keyFilter, _ := cmd.Flags().GetStringArray("key")
		groupFilter, _ := cmd.Flags().GetStringArray("group")

		options := models.CodeTemplateListOptions{
			DatabaseType:    db,
			FilterTypes:     typeFiler,
			FilterLanguages: langFilter,
			FilterModels:    modelFilter,
			FilterKeys:      keyFilter,
			FilterGroups:    groupFilter,
		}
		templates := app.Code.TemplateService.List(&options)
		var tp templatePrinter

		ui.PrintConsoleTitle("ModelHelper Templates")
		fmt.Printf("\nIn the list below you will find all available templates in ModelHelper\n")

		if len(group) > 0 {
			grp := app.Code.TemplateService.Group(group, templates)

			for k, l := range grp {

				// fmt.Println("")
				// fmt.Println("------------------------------------------------------")
				// fmt.Println("")
				ui.PrintConsoleTitle(k)
				tp.templates = l
				ui.RenderTable(&tp)
				fmt.Println("")
			}
		} else {

			tp.templates = templates
			ui.RenderTable(&tp)
		}
	}
}

type templatePrinter struct {
	templates map[string]models.CodeTemplate
}

func (t *templatePrinter) Rows() [][]string {
	var rows [][]string

	for name, t := range t.templates {
		groups := strings.Join(t.Features, ", ")
		row := []string{
			name,
			t.Language,
			t.Type,
			t.Model,
			t.Key,
			groups,
			t.Short,
		}

		rows = append(rows, row)
	}

	return rows
}

func (t *templatePrinter) Header() []string {

	row := []string{
		"Name",
		"Language",
		"Type",
		"Model",
		"Key",
		"Groups",
		"Description",
	}

	return row
}

func openPathInEditor(editor string, loc string) {
	exe := exec.Command(editor, loc)
	if exe.Run() != nil {
		//vim didn't exit with status code 0
	}
}

func getEditor(cfg *models.Config) string {
	if len(cfg.DefaultEditor) > 0 {
		return cfg.DefaultEditor
	} else {
		return ui.PromptForEditor("Please select editor to open the config")
	}
}
