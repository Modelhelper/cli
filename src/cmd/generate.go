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
	"time"

	"modelhelper/cli/app"

	"modelhelper/cli/tpl"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates code based on language, template and source",

	Run: func(cmd *cobra.Command, args []string) {

		charCount := 0

		codeOnly, _ := cmd.Flags().GetBool("code-only")
		isDemo, _ := cmd.Flags().GetBool("demo")
		entities, _ := cmd.Flags().GetStringArray("entity")

		cgf := modelHelperApp.Configuration
		if isDemo == false && len(entities) == 0 {
			return
		}

		//ctx := modelHelperApp.CreateContext()

		inputTemplates, err := cmd.Flags().GetStringArray("template")

		if err != nil {
			panic(err)
		}

		printScreen, _ := cmd.Flags().GetBool("screen")
		start := time.Now()

		var generatedCode []string
		if len(inputTemplates) > 0 {

			tl := tpl.TemplateLoader{
				Directory: app.TemplateFolder(cgf.Templates.Location),
			}

			allTemplates, _ := tl.LoadTemplates()
			//blocks := tpl.ExtractBlocks(&allTemplates)

			for _, tname := range inputTemplates {
				// var tt *tpl.Template
				fmt.Println(tname)
				currentTemplate, found := allTemplates[tname]

				if found {

					if isDemo {
						o, _ := currentTemplate.Generate(testTable())

						generatedCode = append(generatedCode, o)

					} else {

						for _, entity := range entities {
							fmt.Println(entity)
							o, _ := currentTemplate.Generate(testTable())
							generatedCode = append(generatedCode, o)
						}
					}

				}

			}

			if printScreen && len(generatedCode) > 0 {
				screenWriter := tpl.ScreenExporter{}
				for _, s := range generatedCode {
					charCount += len(s)
					screenWriter.Export([]byte(s))
				}
			}

		}

		duration := time.Since(start)

		if !codeOnly {
			con := 1.2
			min := float64(charCount) * con / 60
			fmt.Printf("\nIt took %vms to generate this code. You saved around %v minutes not typing it youreself", duration.Milliseconds(), min)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringArrayP("template", "t", []string{}, "a list of template to convert")
	generateCmd.Flags().StringArrayP("entity", "e", []string{}, "a list of entits to use as a model")
	generateCmd.Flags().Bool("screen", false, "List the output to the screen")
	generateCmd.Flags().String("export", "", "Exports to a directory")
	generateCmd.Flags().Bool("export-bykey", false, "Exports the code using the template keys")
	generateCmd.Flags().Bool("code-only", false, "Writes only the generated code to the console, no stats, no messages - only code")
	generateCmd.Flags().Bool("demo", false, "Uses a demo as input source, this will override any other input sources (entity, graphql) ")
}

func testTable() *tpl.EntityImportModel {
	table := tpl.EntityImportModel{
		Code: tpl.CodeImportModel{
			Language: "cs",
			Creator:  tpl.CreatorImportModel{CompanyName: "Patogen", DeveloperName: "Hans-Petter Eitvet"},
			Types:    testTypes(),
			Imports: []string{
				"using Microsoft.Logging;",
				"using Microsoft.DependencyInjection;",
			},
		},
		Name:           "Contact",
		Description:    "This is a description provided from the table",
		HasDescription: true,
		Columns: []tpl.EntityColumnImportModel{
			{Name: "Id", DataType: "int", IsNullable: false, HasDescription: true, Description: "Description for this column"},
			{Name: "FirstName", DataType: "varchar", IsNullable: false},
			{Name: "LastName", DataType: "varchar", IsNullable: false},
			{Name: "Age", DataType: "int", IsNullable: true},
			{Name: "IsCool", DataType: "bit", IsNullable: true},
			{Name: "TypeId", DataType: "int", IsNullable: true},

			{Name: "ChildTest", DataType: "varchar", IsNullable: true, HasPrefix: true, NameWithoutPrefix: "Test"},
		},
	}
	c1 := tpl.EntityRelation{
		Name: "ContactAddress",
		ReleatedColumn: tpl.EntityColumnProps{
			Name: "ContactId", DataType: "int", IsNullable: true,
		},
		OwnerColumn: tpl.EntityColumnProps{
			Name: "Id", DataType: "int", IsNullable: false,
		},
	}

	p1 := tpl.EntityRelation{
		Name: "ContactType",
		ReleatedColumn: tpl.EntityColumnProps{
			Name: "TypeId", DataType: "int", IsNullable: true,
		},
		OwnerColumn: tpl.EntityColumnProps{
			Name: "Id", DataType: "int", IsNullable: false,
		},
	}

	table.Children = append(table.Children, c1)
	table.Parents = append(table.Parents, p1)
	return &table
}

func testTypes() map[string]tpl.CodeTypeImportModel {
	tl := make(map[string]tpl.CodeTypeImportModel)

	tl["model"] = tpl.CodeTypeImportModel{
		NamePostfix: "",
		NameSpace:   "Testing.Models",
		Key:         "key",
		Imports:     []string{"using HotChocolate;"},
		// Imports:     []string{},
	}
	tl["resolver"] = tpl.CodeTypeImportModel{
		NamePostfix: "Resolver",
		NameSpace:   "Testing.Resolvers",
		Key:         "key",
	}
	tl["inteface"] = tpl.CodeTypeImportModel{
		NamePostfix: "Repository",
		NameSpace:   "Testing.Data",
		NamePrefix:  "I",
		Key:         "key",
	}
	tl["repository"] = tpl.CodeTypeImportModel{
		NamePostfix: "Repository",
		NameSpace:   "Testing.Data",
		Key:         "key",
	}
	return tl
}

// func getEntityModel(name string) interface{} {
// 	src := source

// 	if len(source) == 0 {
// 		src = getSourceName()
// 	}
// 	input := input.GetSource(src, mhConfig)

// 	e, err := input.Entity(name)
// 	if err == nil {
// 		fmt.Println("The entity could not be found")
// 	}

// 	// em := tpl.EntityToModel{
// 	// 	Entity: e,
// 	// }
// 	// m := em.Convert()

// 	return e
// }
