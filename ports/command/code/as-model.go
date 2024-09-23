package code

import (
	"modelhelper/cli/modelhelper"

	"github.com/spf13/cobra"
)

type codeModelOptions struct {
	SourceGroups []string
	Sources      []string
	Templates    []string
	Features     []string
}

type codeModelHandler struct {
	app     *modelhelper.ModelhelperCli
	options *codeModelOptions
}

func NewAsModelCodeCommand(app *modelhelper.ModelhelperCli) *cobra.Command {
	o := &codeModelOptions{}
	h := codeModelHandler{
		app:     app,
		options: o,
	}

	cmd := &cobra.Command{
		Use:   "model",
		Short: "Use this command to see the model for the source or sources",
		Long:  "",
		Run:   h.CommandHandler, // codeCommandHandler(app),
	}

	cmd.Flags().StringArrayVarP(&o.SourceGroups, "source-group", "g", []string{}, "Use a group of source items (must be defined in the current connection)")
	cmd.Flags().StringArrayVarP(&o.Sources, "source", "s", []string{}, "A list of source items to use as a model")
	cmd.Flags().StringArrayVarP(&o.Templates, "template", "t", []string{}, "A list of template to convert")
	cmd.Flags().StringArrayVarP(&o.Features, "feature", "f", []string{}, "Use a group of templates")
	// registerFlags(cmd)

	return cmd
}

func (h *codeModelHandler) CommandHandler(cmd *cobra.Command, args []string) {
	// options := parseCodeOptions(cmd, args)
	// result, err := h.app.Code.Generator.Generate(cmd.Root().Context(), options)
	// if err != nil {
	// 	// handle error
	// 	fmt.Println(err)
	// }
	// sb := strings.Builder{}
}

func SourceModelCommandHandler(app *modelhelper.ModelhelperCli) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		// options := parseCodeOptions(cmd, args)
		// result, err := app.Code.Generator.Generate(cmd.Root().Context(), options)
		// if err != nil {
		// 	// handle error
		// 	fmt.Println(err)
		// }
		// sb := strings.Builder{}

	}

}
