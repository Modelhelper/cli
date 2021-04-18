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
	"fmt"

	"modelhelper/cli/app"
	"modelhelper/cli/tpl"
	"modelhelper/cli/types"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates code based on language, template and source",
	Run: func(cmd *cobra.Command, args []string) {

		// start := time.Now()

		// //generateTestT()
		// duration := time.Since(start)
		// fmt.Printf("\nIt took %vms to generate this code. You saved around 2 hours not typing it youreself", duration.Milliseconds())

		// c := testTemplate()
		// tpl := template.Must(template.New("poco").Parse(c))
		tl := tpl.TemplateLoader{
			Directory: app.TemplateFolder(mhConfig.Templates.Location),
		}

		// ttt, _ := tl.LoadTemplates("", "")
		// if ttt != nil {

		// 	fmt.Println(ttt)
		// }

		tt, _ := tl.LoadTemplate("cs-blocks-poco")
		fmt.Println(tt.Name, tt.Version, tt.Body)
		// tpl.Execute(os.Stdout, entity)

	},
}

func generateTestT() {
	entity := testTable()

	tt := oneMoreTest()

	cnt, err := tt.Generate(entity)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cnt)

}

func init() {
	rootCmd.AddCommand(generateCmd)
}

func testTable() *types.EntityImportModel {
	table := types.EntityImportModel{
		Code: types.CodeImportModel{
			Language: "cs",
			Creator:  types.CreatorImportModel{CompanyName: "Patogen", DeveloperName: "Hans-Petter Eitvet"},
			Types:    testTypes(),
			Imports: []string{
				"using Microsoft.Logging",
				"using Microsoft.DependencyInjection",
			},
		},
		Name: "TestTable",

		NonIgnoredColumns: []types.EntityColumnImportModel{
			{Name: "TestId", DataType: "string"},
			{Name: "FirstName", DataType: "string"},
			{Name: "LastName", DataType: "string"},
		},
	}

	return &table
}

func oneMoreTest() *tpl.Template {
	t := tpl.Template{}

	t.Name = "poco"
	t.Body = testTemplate()

	return &t
}

func testTemplate() string {
	return `
// code her - indent by SPACE not TAB
{{- $model := index .Code.Types "model" }}
using System;
{{- range .Code.Imports }}
{{.}}
{{- end }}

//namespace from block
namespace {{ template "namespace" $model }} {
	public class {{ template "classname" . }} {

	}
}

//namespace from template
namespace {{ $model.NameSpace }} 
{	
	
	public class {{ .Name }}{{ $model.NamePostfix}} 
	{

	}
}

{{.}}
`
}

/*

//namespace from block
namespace {{ template "block1" $model }} {

}

*/
func testTypes() map[string]types.CodeTypeImportModel {
	tl := make(map[string]types.CodeTypeImportModel)

	tl["model"] = types.CodeTypeImportModel{NamePostfix: "Model", NameSpace: "Testing.Test.Test", NamePrefix: "Bla", Key: "key"}
	return tl
}
