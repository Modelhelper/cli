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
	"html/template"
	"modelhelper/cli/types"
	"os"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates code based on language, template and source",
	Run: func(cmd *cobra.Command, args []string) {
		entity := testTable()

		c := testTemplate()
		tpl := template.Must(template.New("poco").Parse(c))

		tpl.Execute(os.Stdout, entity)

	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

func testTemplate() string {
	return `// code her - indent by SPACE not TAB
{{- $model := index .Code.Types "model" }}
using System;
{{- range .Code.Imports }}
{{.}}
{{- end }}

namespace {{ $model.NameSpace}}
{	
	
	public class {{ .Name | Prefix | Pascal | Singular }}{{ $model.NamePrefix }}
	{
	{{- range .NonIgnoredColumns }}
		public {{ .Name }} { get; set; }
	{{- end }}    
	}
}
	
`
}

func testTypes() map[string]types.CodeTypeImportModel {
	tl := make(map[string]types.CodeTypeImportModel)

	tl["model"] = types.CodeTypeImportModel{NamePostfix: "Model", NameSpace: "Testing.Test.Test", NamePrefix: "Bla", Key: "key"}
	return tl
}
