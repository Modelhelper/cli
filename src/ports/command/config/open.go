package config

import (
	"modelhelper/cli/modelhelper/models"
	"modelhelper/cli/ports/config"
	"modelhelper/cli/ui"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewOpenConfigCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "open",
		Short: "Opens the project config file in an editor",
		Long:  "",

		Run: openCommand,
	}

	cmd.Flags().String("editor", "", "The editor to use when opening the file")

	return cmd
}

func openCommand(cmd *cobra.Command, args []string) {
	open, _ := cmd.Flags().GetBool("open")
	ed, _ := cmd.Flags().GetString("editor")

	if open {

		var editor string
		if len(ed) > 0 {
			editor = ed
		} else {
			c := config.NewConfigLoader()
			cfg, err := c.Load()

			if err != nil {
				// handle error
			}
			editor = getEditor(cfg)
		}

		loc := filepath.Join(config.Location(), "config.yaml")
		openPathInEditor(editor, loc)
	}

}

func openPathInEditor(editor string, loc string) {
	exe := exec.Command(editor, loc)
	if exe.Run() != nil {
		//vim didn't exit with status code 0
	}
}

func getEditor(cfg *models.Config) string {
	if len(cfg.DefaultEditor) > 0 {
		return cfg.DefaultEditor
	} else {
		return ui.PromptForEditor("Please select editor to open the config")
	}
}
