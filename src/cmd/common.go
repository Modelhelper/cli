package cmd

import (
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/ui"
	"os/exec"
)

func openPathInEditor(editor string, loc string) {
	exe := exec.Command(editor, loc)
	if exe.Run() != nil {
		//vim didn't exit with status code 0
	}
}

func getEditor(cfg *modelhelper.Config) string {
	if len(cfg.DefaultEditor) > 0 {
		return cfg.DefaultEditor
	} else {
		return ui.PromptForEditor("Please select editor to open the config")
	}
}
