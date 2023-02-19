package config

import "github.com/spf13/cobra"

func NewConfigCommand() *cobra.Command {

	subCommands := []*cobra.Command{
		// NewSetCommand(),
		// NewSetCommand(),
		NewOpenConfigCommand(),
	}

	rootCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage modelhelper configuration",
	}

	for _, sub := range subCommands {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}
