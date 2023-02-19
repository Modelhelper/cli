package source

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ListCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all source items",

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
	fmt.Println("In the source list|ls command")
}
