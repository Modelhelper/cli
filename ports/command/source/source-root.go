package source

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/ui"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func SourceCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	subCommands := []*cobra.Command{
		ListCommand(app),
		NewCopySourceCommand(app),
	}

	rootCmd := &cobra.Command{
		Use:     "source",
		Aliases: []string{"e", "entity", "s", "src"},
		Short:   "Work with items in a source",
		Args:    cobra.RangeArgs(0, 1),
		Run:     rootCommandHandler(app),
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	rootCmd.Flags().StringP("connection", "c", "", "The name of the connection to use")
	rootCmd.Flags().GetString("connection")
	rootCmd.Flags().GetBool("demo")

	return rootCmd
}

func rootCommandHandler(app *modelhelper.ModelhelperCli) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		if len(args) > 0 {

			connection, _ := cmd.Flags().GetString("connection")
			conName, conType := "", ""

			demo, _ := cmd.Flags().GetBool("demo")

			if demo {
				conName = "demo"
				conType = "file"

			} else {

				connections, err := app.ConnectionService.Connections()
				if err != nil {
					// return nil, err
					return
				}
				if len(connections) == 0 {
					return
					// return nil, errors.New("Could not find any connections to use, please add a connection to the config file")
				}
				if len(connection) == 0 {

					conName = app.Config.DefaultConnection
				}

				if len(conName) == 0 {
					for _, v := range connections {
						if strings.Compare(connection, v.Name) == 0 {
							conName = v.Name
							break
						}
					}
				}

				if len(conName) == 0 {
					fmt.Println("Could not find the connection. Please use the --connection flag to specify the connection to use.")

					return
				}
				conType = connections[conName].Type
				// con = g.connectionService.Connection(options.ConnectionName)
			}

			src, _ := app.SourceFactory.NewSource(conType, conName)

			en := args[0]
			e, err := src.Entity(en)
			if err != nil {
				fmt.Println(err)
			}

			if e == nil {
				fmt.Println("The entity could not be found")
				return
			}

			p := message.NewPrinter(language.English)

			// maxL := len(e.Schema) + len(e.Name)

			fmt.Printf("\nEntity:         %s.%s", e.Schema, e.Name)
			fmt.Printf("\nRows:           %v", p.Sprintf("%d", e.RowCount))
			// fmt.Printf("\nIs Versioned:   %s", yesNo(e.IsVersioned))

			if e.IsVersioned {
				fmt.Printf("\nHist. Table:    %s", e.HistoryTable)
			}

			// fmt.Printf("\nCreated: %s\n", "Unknown")

			if len(e.Description) > 0 {
				ui.PrintConsoleTitle("Description:")
				fmt.Println(e.Description)
			}
			renderColumns(&e.Columns)

			if len(e.Indexes) > 0 {
				ui.PrintConsoleTitle("Indexes")

				itr := indexTableRenderer{
					rows: e.Indexes,
				}

				ui.RenderTable(&itr)
			}

			if len(e.ChildRelations) > 0 {
				ui.PrintConsoleTitle("One to many (.ChildRelations)")
				crtr := relTableRenderer{
					rows: e.ChildRelations,
				}

				ui.RenderTable(&crtr)
			}

			if len(e.ParentRelations) > 0 {
				ui.PrintConsoleTitle("Many to one (.ParentRelations)")
				crtr := relTableRenderer{
					rows: e.ParentRelations,
				}
				ui.RenderTable(&crtr)

			}
			fmt.Println("")
		}
	}
}
