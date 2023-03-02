package language

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"modelhelper/cli/ui"

	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func NewListLanguagesCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "List all languages",
		Run:     listlanguagesCommandHandler(app),
	}

	return cmd
}

func listlanguagesCommandHandler(app *modelhelper.ModelhelperCli) func(cmd *cobra.Command, args []string) {

	return func(cmd *cobra.Command, args []string) {

		defs := app.LanguageService.List()

		ui.ConsoleTitle("Language list")
		fmt.Println(`
	This is a list of all available languages defined for model helper			
				`)
		renderer := languageTableRenderer{defs}
		ui.RenderTable(&renderer)
	}
}

type languageTableRenderer struct {
	rows map[string]models.LanguageDefinition
}

func (d *languageTableRenderer) Header() []string {
	h := []string{"Language", "Version", "Datatypes", "Imports", "Keys", "Injects", "Description"}

	return h
}
func (d *languageTableRenderer) Rows() [][]string {
	var rows [][]string

	p := message.NewPrinter(language.English)

	for _, row := range d.rows {
		// un := "No"
		// ci := "No"
		r := []string{
			row.Language,
			row.Version,
			p.Sprintf("%d", len(row.DataTypes)),
			p.Sprintf("%d", len(row.DefaultImports)),
			p.Sprintf("%d", len(row.Keys)),
			p.Sprintf("%d", len(row.Inject)),
			row.Short,
		}

		rows = append(rows, r)
	}

	return rows
}
