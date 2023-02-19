package template

import (
	"fmt"
	"modelhelper/cli/app"
	"modelhelper/cli/config"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/tpl"
	"modelhelper/cli/ui"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func ListCommand() *cobra.Command {

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
		Run: listTemplateCommandHandler,
	}

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

func listTemplateCommandHandler(cmd *cobra.Command, args []string) {
	open, _ := cmd.Flags().GetBool("open")
	cloader := config.NewConfigLoader()
	cfg, _ := cloader.Load()

	tl := tpl.TemplateLoader{
		Directory: app.TemplateFolder(cfg.Templates.Location),
	}
	allTemplates, _ := tl.LoadTemplates()

	if len(args) > 0 {
		current, found := allTemplates[args[0]]
		if found {
			if open {
				editor := getEditor(cfg)
				openPathInEditor(editor, current.TemplateFilePath)
			} else {

				fmt.Print(current.ToString(args[0]))
			}
		}

		return
	}
	group, _ := cmd.Flags().GetString("by")
	typeFiler, _ := cmd.Flags().GetStringArray("type")
	langFilter, _ := cmd.Flags().GetStringArray("lang")
	modelFilter, _ := cmd.Flags().GetStringArray("model")
	keyFilter, _ := cmd.Flags().GetStringArray("key")
	groupFilter, _ := cmd.Flags().GetStringArray("group")

	var tp templatePrinter
	if len(typeFiler) > 0 {
		ft := tpl.FilterByType{}
		allTemplates = ft.Filter(allTemplates, typeFiler)
	}

	if len(langFilter) > 0 {
		ft := tpl.FilterByLang{}
		allTemplates = ft.Filter(allTemplates, langFilter)
	}
	if len(keyFilter) > 0 {
		ft := tpl.FilterByKey{}
		allTemplates = ft.Filter(allTemplates, keyFilter)
	}
	if len(modelFilter) > 0 {
		ft := tpl.FilterByModel{}
		allTemplates = ft.Filter(allTemplates, modelFilter)
	}
	if len(groupFilter) > 0 {
		ft := tpl.FilterByGroup{}
		allTemplates = ft.Filter(allTemplates, groupFilter)
	}

	if len(group) > 0 {
		ui.PrintConsoleTitle("ModelHelper Templates grouped by " + group)
		fmt.Printf("\nIn the list below you will find all available templates in ModelHelper\n")

		grouper := tpl.GetGrouper(strings.ToLower(group))
		descr := tpl.GetDescriber(group)

		mg := grouper.Group(allTemplates)

		for typ, tv := range mg {
			ui.PrintConsoleTitle(typ)

			desc := descr.Describe(typ)
			if desc != nil {
				fmt.Println(desc.Long)
			}

			fmt.Println()
			tp.templates = tv
			ui.RenderTable(&tp)

		}
	} else {
		ui.PrintConsoleTitle("ModelHelper Templates")
		fmt.Printf("\nIn the list below you will find all available templates in ModelHelper\n")

		tp.templates = allTemplates
		ui.RenderTable(&tp)
	}

	fmt.Println()
}

type templatePrinter struct {
	templates map[string]tpl.Template
}

func (t *templatePrinter) Rows() [][]string {
	var rows [][]string

	for name, t := range t.templates {
		groups := strings.Join(t.Groups, ", ")
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

func getEditor(cfg *modelhelper.Config) string {
	if len(cfg.DefaultEditor) > 0 {
		return cfg.DefaultEditor
	} else {
		return ui.PromptForEditor("Please select editor to open the config")
	}
}
