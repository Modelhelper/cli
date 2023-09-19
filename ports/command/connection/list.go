package connection

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"modelhelper/cli/ui"

	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func NewListConnectionsCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "List all connections",
		Run:     listConnectionCommandHandler(app),
	}

	return cmd
}

func listConnectionCommandHandler(app *modelhelper.ModelhelperCli) func(cmd *cobra.Command, args []string) {

	return func(cmd *cobra.Command, args []string) {

		connections, err := app.ConnectionService.Connections()

		if err != nil {
			fmt.Errorf("Err when listing connections %w", err)
		}
		ui.PrintConsoleTitle("Connections list")
		fmt.Println(`
This is a list of all available connections you can use as a source input for templates			
				`)
		renderer := connectionTableRenderer{connections}
		ui.RenderTable(&renderer)
	}
}

type connectionTableRenderer struct {
	rows map[string]*models.ConnectionList
}

func (d *connectionTableRenderer) Header() []string {
	h := []string{"Name", "Type", "Default", "Groups", "Synonyms", "Options", "Description"}

	return h
}
func (d *connectionTableRenderer) Rows() [][]string {
	var rows [][]string

	p := message.NewPrinter(language.English)

	for _, row := range d.rows {
		// un := "No"
		// ci := "No"
		id := "No"

		if row.IsDefault {
			id = "Yes"
		}
		r := []string{
			row.Name,
			row.Type,
			id,
			p.Sprintf("%d", len(row.Groups)),
			p.Sprintf("%d", len(row.Synonyms)),
			p.Sprintf("%d", len(row.Options)),
			row.Description,
		}

		rows = append(rows, r)
	}

	return rows
}
