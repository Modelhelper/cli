package code

import (
	"modelhelper/cli/modelhelper"

	"github.com/spf13/cobra"
)

func NewCodeRootCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	subCommands := []*cobra.Command{
		NewGenerateCodeCommand(app),
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
