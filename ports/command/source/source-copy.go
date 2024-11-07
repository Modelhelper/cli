package source

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/ports/source/mssql"

	"github.com/spf13/cobra"
)

type copySourceHandler struct {
	connectionService modelhelper.ConnectionService
	sourceFactory     modelhelper.SourceFactoryService
}

func NewCopySourceCommand(app *modelhelper.ModelhelperCli) *cobra.Command {
	h := copySourceHandler{app.ConnectionService, app.SourceFactory}
	cmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy data from one source to another",

		Run: h.copyCommandHandler,
	}

	cmd.Flags().StringP("source", "s", "", "The name of the connection to use as source")
	cmd.Flags().StringP("dest", "d", "", "The name of the connection to use as destination")
	cmd.Flags().StringP("query", "q", "", "The query to use to get the data")
	cmd.Flags().StringP("table", "t", "", "The table to use to get the data")

	return cmd
}

func (h copySourceHandler) copyCommandHandler(cmd *cobra.Command, args []string) {

	bc := mssql.NewMsSqlDatabaseService(h.connectionService, h.sourceFactory)

	src, _ := cmd.Flags().GetString("source")
	dest, _ := cmd.Flags().GetString("dest")
	query, _ := cmd.Flags().GetString("query")
	table, _ := cmd.Flags().GetString("table")

	if query == "" {
		query = "SELECT * FROM " + table
	}

	rows, err := bc.BulkCopy(src, dest, query, table)
	if err != nil {
		cmd.Println(err)
	}

	fmt.Printf("Finished copying %d rows\n", rows)
	// cmd.Println(rows)
}
