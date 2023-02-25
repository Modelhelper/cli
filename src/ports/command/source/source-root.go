package source

import (
	"fmt"

	"github.com/spf13/cobra"
)

func SourceCommand() *cobra.Command {

	subCommands := []*cobra.Command{
		ListCommand(),
	}

	rootCmd := &cobra.Command{
		Use:     "source",
		Aliases: []string{"e", "entity", "s", "src"},
		Short:   "Work with items in a source",
		Args:    cobra.RangeArgs(0, 1),
		Run:     rootCommandHandler,
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}

func rootCommandHandler(cmd *cobra.Command, args []string) {
	fmt.Println("source command root")

	if len(args) > 0 {
		fmt.Println("Read details")
	}
}
