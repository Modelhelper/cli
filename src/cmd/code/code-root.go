package code

import "github.com/spf13/cobra"

func NewGenerateCommand() *cobra.Command {

	subCommands := []*cobra.Command{
		NewGenerateCodeCommand(),
	}

	rootCmd := &cobra.Command{
		Use:     "code",
		Aliases: []string{"c"},
		Short:   "Generates code based on language, template and source",
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}
