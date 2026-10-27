package command

import (
	"github.com/Modelhelper/cli/v3/modelhelper/models"
	"github.com/Modelhelper/cli/v3/ui"
	"os/exec"
)

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
