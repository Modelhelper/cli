package template

import (
	"modelhelper/cli/app"
	"modelhelper/cli/config"
	"modelhelper/cli/tpl"
	"modelhelper/cli/ui"

	"github.com/spf13/cobra"
)

func OpenCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "open",
		Short: "Opens the template file in an editor",
		Long:  "",
		Args:  cobra.RangeArgs(1, 2),

		Run: openCommand,
	}

	cmd.Flags().String("editor", "", "The editor to use when opening the file")

	return cmd
}

func openCommand(cmd *cobra.Command, args []string) {
	cloader := config.NewConfigLoader()
	cfg, _ := cloader.Load()
	editor := getEditor(cfg)

	tl := tpl.TemplateLoader{
		Directory: app.TemplateFolder(cfg.Templates.Location),
	}
	allTemplates, _ := tl.LoadTemplates()

	if len(args) > 0 {

		if len(args) == 2 {
			editor = args[1]
		}

		if len(editor) == 0 {
			ui.PromptForEditor("Select editor")
		}

		current, found := allTemplates[args[0]]
		if found {
			openPathInEditor(editor, current.TemplateFilePath)
		}

	}
}
