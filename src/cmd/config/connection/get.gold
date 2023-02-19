package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewGetCommand() *cobra.Command {

	generateCmd := &cobra.Command{
		Use:   "get",
		Short: "Gets a config value",
		Long:  "",
		Run:   getCommandHandler,
	}

	generateCmd.Flags().String("scope", "root", "The scope must be either root or project")

	generateCmd.Flags().Bool("screen", false, "List the output to the screen, default false")
	generateCmd.Flags().Bool("copy", false, "Copies the generated code to the clipboard (ctrl + v), default false")
	generateCmd.Flags().String("export-path", "", "Exports to a directory")
	// generateCmd.Flags().Bool("export-bykey", false, "Exports the code using the template keys, default false")
	generateCmd.Flags().Bool("overwrite", false, "Overwrite any existing file when exporting to file on disk")

	generateCmd.Flags().Bool("code-only", false, "Writes only the generated code to the console, no stats, no messages - only code, default false")

	generateCmd.Flags().Bool("demo", false, "Uses a demo as input source, this will override any other input sources (entity, graphql), default false ")

	generateCmd.Flags().String("config-path", "", "Instructs the program to use this config as the config")
	generateCmd.Flags().String("project-path", "", "Instructs the program to use this project as input")

	generateCmd.Flags().String("key", "", "The key to use when encoding and decoding secrets for a connection")

	// generateCmd.Flags().String("setup", "", "Use this setup to generate code") // version 3.1
	generateCmd.Flags().StringP("connection", "c", "", "The connection key to be used, uses default connection if not provided")

	return generateCmd
}

func getCommandHandler(cmd *cobra.Command, args []string) {
	fmt.Println("The config set command")
}
