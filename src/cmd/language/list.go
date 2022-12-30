package language

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewListLanguagesCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "List all languages",
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
		Run: listlanguagesCommandHandler,
	}

	cmd.Flags().String("editor", "", "The editor to use when opening the file")

	return cmd
}

func listlanguagesCommandHandler(cmd *cobra.Command, args []string) {
	fmt.Println("language list")
}
