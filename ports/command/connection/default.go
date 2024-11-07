package connection

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"

	"github.com/spf13/cobra"
)

type defaultConnectionHandler struct {
	connectionService modelhelper.ConnectionService
	configService     modelhelper.ConfigService
	config            *models.Config
	// program           *tea.Program
}

type selectItem struct {
	Name        string
	Description string
}

var selectedItem selectItem

func NewSetDefaultConnectionCommand(cs modelhelper.ConnectionService, cfg *models.Config, cfgSrv modelhelper.ConfigService) *cobra.Command {
	handler := &defaultConnectionHandler{
		connectionService: cs,
		configService:     cfgSrv,
		config:            cfg,
	}

	cmd := &cobra.Command{
		Example: "mh connection default 'name'",
		Use:     "default [name]",
		Args:    cobra.MaximumNArgs(1),

		Short: "Sets a new default connection",
		RunE:  handler.handleCmd,
	}
	return cmd
}

func (h *defaultConnectionHandler) handleCmd(cmd *cobra.Command, args []string) error {
	con := ""

	if len(args) > 0 {
		con = args[0]
	}

	if con != "" && h.config.DefaultConnection != con {

		h.config.DefaultConnection = con
		err := h.configService.SaveConfig(h.config)

		if err != nil {
			fmt.Printf("Failed to set default connection: %s", err)
			return err
		}

		fmt.Printf("'%s' is the new default connection", con)

		return nil
	}

	con, err := getConnectionNameFromSelectionList(h.connectionService)
	// con is blank, let user select from list
	// cons, err := h.connectionService.Connections()
	if err != nil {
		return err
	}
	// maxLen := 0

	// keys := make([]string, 0, len(cons))
	// for k := range cons {
	// 	if len(k) > maxLen {
	// 		maxLen = len(k)
	// 	}
	// 	keys = append(keys, k)
	// }

	// sort.Strings(keys)
	// opts := []huh.Option[string]{}

	// for _, k := range keys {
	// 	o := cons[k]
	// 	ho := huh.NewOption(fmt.Sprintf("%-*s [%-*s] - %s", maxLen, k, 8, o.Type, o.Description), k)
	// 	opts = append(opts, ho)
	// }

	// form := huh.NewForm(
	// 	huh.NewGroup(
	// 		huh.NewSelect[string]().
	// 			Title("Select a connection from the list").
	// 			Height(10).
	// 			Value(&con).
	// 			Options(opts...),
	// 	),
	// )

	// err = form.Run()
	// if err != nil {
	// 	return err
	// }

	if con == "" {
		return nil
	}

	h.config.DefaultConnection = con
	err = h.configService.SaveConfig(h.config)

	if err != nil {
		return err
	}

	fmt.Printf("'%s' is the new default connection", selectedItem.Name)

	return nil
}
