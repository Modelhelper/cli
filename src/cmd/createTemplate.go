/*
Copyright © 2020 Hans-Petter Eitvet

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"io/ioutil"
	"modelhelper/cli/defaults"

	"github.com/spf13/cobra"
)

var createTplLang string = "cs"
var createTplPath string = "cs"
var createTplKey string = "cs"
var createTplImport string = "table"

// createTemplateCmd represents the createTemplate command
var createTemplateCmd = &cobra.Command{
	Use:     "template",
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

func init() {
	createCmd.AddCommand(createTemplateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createTemplateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createTemplateCmd.Flags().StringVarP(&createTplLang, "language", "l", "cs", "Defines the template language")
	createTemplateCmd.Flags().StringVarP(&createTplPath, "path", "p", "c:\\temp\\templates", "Defines the template language")
	createTemplateCmd.Flags().StringVarP(&createTplImport, "import", "i", "cs", "Defines the template language")
	createTemplateCmd.Flags().StringVarP(&createTplKey, "key", "k", "cs", "Defines the template language")
}
