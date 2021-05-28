package cmd

import (
	"fmt"
	"modelhelper/cli/code"
	"modelhelper/cli/config"

	"github.com/spf13/cobra"
)

// aboutCmd represents the about command
var languageCmd = &cobra.Command{
	Use:     "language",
	Aliases: []string{"lang", "l"},
	Short:   "Root command for working with language definitions installed for modelhelper",
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()

		defs, err := code.LoadFromPath(cfg.Languages.Definitions)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		if len(args) > 0 {
			def := defs[args[0]]
			fmt.Println(def)
			// showLanguage(args[0], cfg.Languages.Definitions)
		} else {
			for k, v := range defs {
				fmt.Print(k, v.Version)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(languageCmd)
	languageCmd.Flags().Bool("open", false, "Opens the language definition in the default editor")
}
