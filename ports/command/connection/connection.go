package connection

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"sort"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func NewConnectionCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	subCommands := []*cobra.Command{
		NewListConnectionsCommand(app),
		NewCreateConnectionCommand(app.ConnectionService),
		NewDeleteConnectionCommand(app.ConnectionService),
		NewSetDefaultConnectionCommand(app.ConnectionService, app.Config, app.ConfigService),
	}

	rootCmd := &cobra.Command{
		Use:   "connection",
		Short: "Manage modelhelper connections",
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}

func getConnectionNameFromSelectionList(cs modelhelper.ConnectionService) (string, error) {
	var connectionName string
	cons, err := cs.Connections()
	if err != nil {
		return "", err
	}
	maxLen := 0

	keys := make([]string, 0, len(cons))
	for k := range cons {
		if len(k) > maxLen {
			maxLen = len(k)
		}
		keys = append(keys, k)
	}

	sort.Strings(keys)
	opts := []huh.Option[string]{}

	for _, k := range keys {
		o := cons[k]
		ho := huh.NewOption(fmt.Sprintf("%-*s [%-*s] - %s", maxLen, k, 8, o.Type, o.Description), k)
		opts = append(opts, ho)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a connection from the list").
				Height(10).
				Value(&connectionName).
				Options(opts...),
		),
	)

	err = form.Run()
	if err != nil {
		return "", err
	}

	return connectionName, nil
}
