package project

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewProjectCommand() *cobra.Command {

	subCommands := []*cobra.Command{
		NewGenerateProjectCommand(),
		NewOpenProjectCommand(),
		NewProjectInitCommand(),
	}

	rootCmd := &cobra.Command{
		Use:     "project",
		Aliases: []string{"p"},
		Short:   "Generates code based on language, template and source",
		Run:     rootCommandHandler,
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}

func rootCommandHandler(cmd *cobra.Command, args []string) {
	fmt.Println("project info")
}
