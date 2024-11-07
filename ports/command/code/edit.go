package code

import (
	"fmt"
	"log/slog"
	"modelhelper/cli/modelhelper"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type codeEditOptions struct {
	SourceGroups []string
	Sources      []string
	Templates    []string
	Features     []string
}

type codeEditHandler struct {
	app     *modelhelper.ModelhelperCli
	options *codeModelOptions
}

func NewEditCodeTemplateCommand(app *modelhelper.ModelhelperCli) *cobra.Command {
	o := &codeModelOptions{}
	h := codeEditHandler{
		app:     app,
		options: o,
	}

	cmd := &cobra.Command{
		Use: "edit",

		Short: "Use this command to open the code template in an editor",
		Long:  "",
		Args:  cobra.RangeArgs(1, 2),

		Run: h.CommandHandler,
	}

	// cmd.Flags().String("editor", "", "Clone the connection from the given name")

	return cmd
}

func (h *codeEditHandler) CommandHandler(cmder *cobra.Command, args []string) {
	editor := h.app.Config.DefaultEditor

	if len(args) > 1 {
		editor = args[1]
	}

	if len(editor) == 0 {
		editor = "nano"
	}

	tpl := h.app.Code.TemplateService.Load(args[0])

	if tpl == nil {
		slog.Info("Template not found", "name", args[0])
		return
	}

	exe := exec.Command(editor, tpl.TemplateFilePath)

	if editor != "code" {
		exe.Stdout = os.Stdout
		exe.Stdin = os.Stdin
		exe.Stderr = os.Stderr
	}

	err := exe.Run()
	if err != nil {
		fmt.Print(err)
		slog.Error("Failed to start editor")
	}
	// if exe.Run() != nil {
	// 	//vim didn't exit with status code 0
	// }
}
