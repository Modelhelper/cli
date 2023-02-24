package template

import (
	"io/ioutil"
	"modelhelper/cli/defaults"
	"modelhelper/cli/modelhelper"

	"github.com/spf13/cobra"
)

var createTplLang string = "cs"
var createTplPath string = "cs"
var createTplKey string = "cs"
var createTplImport string = "table"

// createTemplateCmd represents the createTemplate command
func CreateCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	var createTemplateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"t", "tpl"},
		Short:   "Creates a template file for a specific language",

		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			t := defaults.TemplateOptions{
				Import:      createTplImport,
				LocationKey: createTplKey,
				Language:    createTplLang,
			}

			b := defaults.DefaultTemplateContent(&t)

			// f, err := os.Create(createTplPath + "/" + name + ".yaml")

			err := ioutil.WriteFile(createTplPath+"\\"+createTplLang+"\\"+name+".yaml", b, 0644)
			if err != nil {
				panic(err)
			}
			// defer f.Close()

			// n, err := f.Write(b)
			// if err != nil {
			// 	panic(err)
			// }

			// fmt.Printf("wrote %d bytes\n", n)
		},
	}

	createTemplateCmd.Flags().StringVarP(&createTplLang, "language", "l", "cs", "Defines the template language")
	createTemplateCmd.Flags().StringVarP(&createTplPath, "path", "p", "c:\\temp\\templates", "Defines the template language")
	createTemplateCmd.Flags().StringVarP(&createTplImport, "import", "i", "cs", "Defines the template language")
	createTemplateCmd.Flags().StringVarP(&createTplKey, "key", "k", "cs", "Defines the template language")

	return createTemplateCmd
}
