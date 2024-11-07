package code

import (
	"fmt"
	"modelhelper/cli/modelhelper"

	"github.com/spf13/cobra"
)

func NewTemplatesCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "template [command] [options]",
		Aliases: []string{"t"},
		Short:   "Work with code templates",

		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("mh project template")
		},
	}

	cmd.AddCommand(NewTemplateListCommand(app))
	cmd.AddCommand(templateDetailCommand(app))
	cmd.AddCommand(NewEditCodeTemplateCommand(app))

	return cmd
}

func templateDetailCommand(app *modelhelper.ModelhelperCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detail [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Get details about a named template",

		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("mh project template detail")
			tpl := app.Project.TemplateService.Load(args[0])

			fmt.Printf("name: %s, lang: %s", tpl.Name, tpl.Language)
		},
	}

	return cmd
}

// func templateDetailCommand(app *modelhelper.ModelhelperCli) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "detail [name]",
// 		Args:  cobra.ExactArgs(1),
// 		Short: "Get details about a named template",

// 		Run: func(cmd *cobra.Command, args []string) {

// 			fmt.Println("mh project template detail")
// 			tpl := app.Project.TemplateService.Load(args[0])

// 			fmt.Printf("name: %s, lang: %s", tpl.Name, tpl.Language)
// 		},
// 	}

// 	return cmd
// }
