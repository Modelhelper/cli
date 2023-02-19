package template

import "github.com/spf13/cobra"

func NewTemplateCommand() *cobra.Command {

	subCommands := []*cobra.Command{
		ListCommand(),
		OpenCommand(),
		CreateCommand(),
	}

	rootCmd := &cobra.Command{
		Use:   "template",
		Short: "Manage modelhelper templates",
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}
