package code

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"sort"

	"modelhelper/cli/modelhelper/models"
	"modelhelper/cli/ui"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/spf13/cobra"
)

type presentSimulationListModel struct {
	templates []models.CodeTemplate
	options   models.CodeTemplateListOptions
}

func NewTemplateListCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

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

	cmd.Flags().StringArray("hide", []string{}, "Hide the columns in the table")

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
		hideCols, _ := cmd.Flags().GetStringArray("hide")

		mc := make(map[string]bool)

		for _, c := range hideCols {
			mc[c] = true
		}

		options := models.CodeTemplateListOptions{
			DatabaseType:    db,
			FilterTypes:     typeFiler,
			FilterLanguages: langFilter,
			FilterModels:    modelFilter,
			FilterKeys:      keyFilter,
			FilterGroups:    groupFilter,
			HideColumns:     mc,
		}
		templates := app.Code.TemplateService.List(&options)
		tp := &templatePrinter{options: options}

		ui.PrintConsoleTitle("ModelHelper Templates")
		fmt.Printf("\nIn the list below you will find all available templates in ModelHelper\n")

		if len(group) > 0 {
			grp := app.Code.TemplateService.Group(group, templates)

			for k, l := range grp {

				// fmt.Println("")
				// fmt.Println("------------------------------------------------------")
				// fmt.Println("")
				ui.PrintConsoleTitle(k)
				tp.templates = templatesByName(l)
				ui.RenderTable(tp)
				fmt.Println("")
			}
		} else {

			tp.templates = templatesByName(templates)
			ui.RenderTable(tp)
		}
	}
}

type templatePrinter struct {
	templates []models.CodeTemplate
	options   models.CodeTemplateListOptions
}

func (tp *templatePrinter) Rows() [][]string {
	var rows [][]string

	for _, t := range tp.templates {
		ok := false
		groups := strings.Join(t.Features, ", ")
		row := []string{
			t.Name,
			t.Language,

			// t.Type,
			// t.Model,
			// t.Key,
			// groups,
			// t.Short,
		}

		_, ok = tp.options.HideColumns["type"]
		if !ok {
			row = append(row, t.Type)
		}
		_, ok = tp.options.HideColumns["model"]
		if !ok {
			row = append(row, t.Model)
		}
		_, ok = tp.options.HideColumns["key"]
		if !ok {
			row = append(row, t.Key)
		}

		_, ok = tp.options.HideColumns["groups"]
		if !ok {

			row = append(row, groups)
		}
		_, ok = tp.options.HideColumns["desc"]
		if !ok {
			row = append(row, t.Short)
		}

		// row = append(row, t.Short)

		rows = append(rows, row)
	}

	return rows
}

func (tp *templatePrinter) Header() []string {
	ok := false
	row := []string{
		"Name",
		"Language",
		// "Type",
		// "Model",
		// "Key",
		// "Groups",
		// "Description",
	}

	_, ok = tp.options.HideColumns["type"]
	if !ok {
		row = append(row, "Type")
	}
	_, ok = tp.options.HideColumns["model"]
	if !ok {
		row = append(row, "Model")
	}
	_, ok = tp.options.HideColumns["key"]
	if !ok {
		row = append(row, "Key")
	}

	_, ok = tp.options.HideColumns["groups"]
	if !ok {

		row = append(row, "Groups")
	}
	_, ok = tp.options.HideColumns["desc"]
	if !ok {
		row = append(row, "Description")
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

func templatesByName(templates map[string]models.CodeTemplate) []models.CodeTemplate {
	var templateArray []models.CodeTemplate

	for _, template := range templates {
		templateArray = append(templateArray, template)
	}

	sort.Slice(templateArray, func(i, j int) bool {
		return templateArray[i].Name < templateArray[j].Name
	})
	return templateArray
}

func (m presentSimulationListModel) Init() tea.Cmd {
	return nil
}

func (m presentSimulationListModel) getPagedata() {
}

func (m presentSimulationListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	// case errMsg:
	// 	m.err = msg
	// 	return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			// m.quitting = true
			return m, tea.Quit
			// case "h", "left":
			// 	fmt.Printf("\n\nprevious\n")
			// case "r", "right":
			// 	fmt.Printf("\n\next\n")
			// // if m.paginator.HasPreviousPage() {
			// // 	m.paginator.PreviousPage()

			// // 	page, err := m.app.Query.Antulator.GetSimulationList.Handle(m.ctx, m.pageSize, m.paginator.NextPage() .Page())
			// // 	if err != nil {
			// // 		m.err = err
			// // 		return m, nil
			// // 	}

			// // 	m.currentPage = page
			// // 	m.simulations = page.Simulations
			// // }
			// case "enter":
			// 	if m.currentPage != nil && m.currentPage.PageNumber < m.currentPage.TotalPages {

			// 		nextPage := m.currentPage.PageNumber + 1

			// 		page, err := m.app.Query.Antulator.GetSimulationList.Handle(m.ctx, m.pageSize, nextPage)

			// 		if err != nil {
			// 			m.err = err
			// 			fmt.Printf("\n\nError: %s\n", err)
			// 			return m, nil
			// 		}

			// 		m.currentPage = page

			// 		// m.simulations = page.Simulations
			// 		m.simulations = append(m.simulations, page.Simulations...)
			// 	}

			// 	if m.currentPage != nil && m.currentPage.PageNumber == m.currentPage.TotalPages {
			// 		fmt.Printf("\nNo more pages\n")
			// 		return m, tea.Quit
			// 	}
			// 	return m, nil

			// default:
			// 	return m, nil
		}

	default:
		return m, nil
	}

	return m, cmd
}

func (m presentSimulationListModel) View() string {
	var b strings.Builder

	// if m.err != nil {
	// 	return m.err.Error()
	// }
	b.WriteString(fmt.Sprintf("\n\n\tTemplates\n\t-------------------\n\n"))
	var rows []table.Row
	for _, template := range m.templates {
		rows = append(rows, table.Row{
			template.Name,
			template.Language,
			template.Type,
			template.Model,
			template.Key,
			template.Features,
			template.Description,
		})
	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"Name", "Language", "Name", "Duration", "Wasp", "Mac", "Use Room sensor", "Bypass"})
	t.AppendRows(rows)
	table := t.Render()

	b.WriteString(table)

	// if m.currentPage != nil && m.currentPage.PageNumber < m.currentPage.TotalPages {
	// 	b.WriteString(fmt.Sprintf("\n\nPage:%d of %d\n", m.currentPage.PageNumber, m.currentPage.TotalPages))
	// 	b.WriteString("\nPress enter to see more. Press q, esc or ctr+c to quit.")
	// }

	// if m.quitting {
	// 	b.WriteString("\n")
	// }

	// b.WriteString("  " + m.paginator.View())
	// b.WriteString("\n\n  h/l ←/→ page • q: quit\n")
	return b.String()

	// return str
}
