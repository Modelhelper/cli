package code

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/ports/extractor"

	"github.com/spf13/cobra"
)

func NewCodeRootCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	subCommands := []*cobra.Command{
		NewGenerateCodeCommand(app),
		NewChangelogCommand(app),
		NewVariableExtractorCommand(),
		NewTemplatesCommand(app),
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

func NewVariableExtractorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "extract",
		Aliases: []string{"e"},
		Short:   "Extract variables",
		Run: func(cmd *cobra.Command, args []string) {
			fn, _ := cmd.Flags().GetString("file")

			if len(fn) == 0 {
				return
			}

			e := extractor.NewDockerStackVariableExtractor(fn)
			list := e.Extract()

			for _, v := range list {

				fmt.Println(v)
			}
		},
	}

	cmd.Flags().StringP("file", "f", "", "The file to extract from")

	return cmd
}
