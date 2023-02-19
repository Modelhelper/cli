package language

import "github.com/spf13/cobra"

func LanguageCommand() *cobra.Command {

	subCommands := []*cobra.Command{
		NewListLanguagesCommand(),
	}

	rootCmd := &cobra.Command{
		Use:   "language",
		Short: "Manage modelhelper languages",
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}
