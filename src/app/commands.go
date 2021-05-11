package app

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type CommandItem struct {
	Name           string
	Short          string
	HasSubCommands bool
	Aliases        string
}

func CommandInfo(cmd *cobra.Command) string {

	cl, err := Commands(cmd)
	if err != nil {
		fmt.Printf("Error when fetching commands: %e", err)
	}

	var sb strings.Builder

	for _, ci := range cl {

		name := ci.Name

		if len(ci.Aliases) > 0 {
			name = ci.Name + " [" + ci.Aliases + "]"
		}

		fmt.Fprintf(&sb, "  %s:%s%s\n", name, " ", ci.Short)
	}

	return fmt.Sprintf(`
  Commands
  ------------

%v
`, sb.String())
}

// Commands returns a list of available commands
func Commands(cmd *cobra.Command) ([]CommandItem, error) {

	list := []CommandItem{}

	for _, c := range cmd.Commands() {

		if c.Hidden {
			continue
		}
		alias, sep := "", ""

		if len(c.Aliases) > 0 {
			for i, a := range c.Aliases {
				if i > 0 {
					sep = ", "
				}
				alias += sep + a
			}
		}

		item := CommandItem{
			c.Name(),
			c.Short,
			true,
			alias,
		}

		list = append(list, item)
	}

	return list, nil
}
