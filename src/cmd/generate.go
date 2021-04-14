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
	tpl "modelhelper/cli/tpl"
	"modelhelper/cli/types"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates code based on language, template and source",
	Run: func(cmd *cobra.Command, args []string) {
		entity := testTable()

		tt := oneMoreTest()

		cnt, err := tt.Generate(entity)

		if err != nil {

		}

		fmt.Println(cnt)

		// c := testTemplate()
		// tpl := template.Must(template.New("poco").Parse(c))

		// tpl.Execute(os.Stdout, entity)

	},
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

	t.Language = "poco"
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
`
}

func testTypes() map[string]types.CodeTypeImportModel {
	tl := make(map[string]types.CodeTypeImportModel)

	tl["model"] = types.CodeTypeImportModel{NamePostfix: "Model", NameSpace: "Testing.Test.Test", NamePrefix: "Bla", Key: "key"}
	return tl
}
