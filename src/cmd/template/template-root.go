package template

import "github.com/spf13/cobra"

func NewTemplateCommand() *cobra.Command {

	subCommands := []*cobra.Command{
		NewListProjectsCommand(),
		NewOpenTemplateCommand(),
	}

	rootCmd := &cobra.Command{
		Use:   "template",
		Short: "Manage modelhelper configuration",
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}
