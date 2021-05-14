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
	"log"
	"modelhelper/cli/app"
	"modelhelper/cli/codegen"
	"modelhelper/cli/config"
	"modelhelper/cli/model"
	"modelhelper/cli/project"
	"modelhelper/cli/source"
	"modelhelper/cli/tpl"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g", "gen"},
	Short:   "Generates code based on language, template and source",

	Run: func(cmd *cobra.Command, args []string) {
		codeOnly, _ := cmd.Flags().GetBool("code-only")
		isDemo, _ := cmd.Flags().GetBool("demo")
		entityFlagArray, _ := cmd.Flags().GetStringArray("entity")
		tempPath, _ := cmd.Flags().GetString("template-path")
		projectPath, _ := cmd.Flags().GetString("project")

		// if isDemo == false && len(entityFlagArray) == 0 {
		// 	return
		// }
		ctx := modelHelperApp.CreateContext()
		// var ctx *app.Context
		var cfg *config.Config
		var prj *project.Project
		var entities []source.Entity

		charCount := 0

		if isDemo {
			cfg = config.Load()
			// load demo project
			var ds *source.DemoSource
			el, _ := ds.Entities("*")

			for _, eitem := range *el {
				entities = append(entities, eitem)
			}

			// load demo tables (2)
		} else {

			conName := ctx.DefaultConnection

			con := ctx.Connections[conName]
			src := con.LoadSource()

			if len(entityFlagArray) > 0 {
				for _, entityName := range entityFlagArray {
					entity, err := src.Entity(entityName)
					if err != nil {
						log.Fatalln(err)
					}

					entities = append(entities, *entity)
				}
			}

			if len(tempPath) > 0 {
				fmt.Println("Use this path to template")
			}

			configFile, _ := cmd.Flags().GetString("config")

			if len(configFile) > 0 {
				fmt.Println("Use this as configuration file")
				cfg = config.LoadFromFile(configFile)
			} else {
				cfg = config.Load()
			}

			// currentProject := *project.Project{}
			if len(projectPath) > 0 {
				prj, _ = project.Load(projectPath)
			} else {
				prj, _ = project.Load(project.DefaultLocation())
			}
		}

		tl := tpl.TemplateLoader{
			Directory: app.TemplateFolder(cfg.Templates.Location),
		}

		allTemplates, _ := tl.LoadTemplates()
		blocks := tpl.ExtractBlocks(&allTemplates)

		inputTemplates, err := cmd.Flags().GetStringArray("template")

		if err != nil {
			panic(err)
		}

		printScreen, _ := cmd.Flags().GetBool("screen")
		start := time.Now()

		var generatedCode []string
		if len(inputTemplates) > 0 {

			for _, tname := range inputTemplates {

				// var tt *tpl.Template
				// fmt.Println(tname)
				currentTemplate, found := allTemplates[tname]

				if found {
					var input interface{}

					tplMap := make(map[string]string)

					for k, b := range blocks {
						tplMap[k] = b.Body
					}
					tplMap[tname] = currentTemplate.Body

					generator := codegen.GoLangGenerator{
						Templates:    tplMap,
						TemplateName: tname,
					}

					if len(currentTemplate.Model) == 0 || currentTemplate.Model == "basic" {
						basicModel := basicModel{
							project: prj,
							key:     currentTemplate.Key,
						}
						input = basicModel.ToModel()
						o, _ := generator.Generate(input)
						generatedCode = append(generatedCode, o)

					}

					if currentTemplate.Model == "entity" && len(entities) > 0 {
						for _, entity := range entities {
							entityModel := entityModel{
								entity:  &entity,
								project: prj,
								key:     currentTemplate.Key,
							}
							input = entityModel.ToModel()
							// fmt.Println(input)
							o, _ := generator.Generate(input)
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

	generateCmd.Flags().String("template-path", "", "Instructs the program to use this path as root for templates")
	generateCmd.Flags().String("config", "", "Instructs the program to use this config as the config")
	generateCmd.Flags().String("project", "", "Instructs the program to use this project as input")

	generateCmd.Flags().String("setup", "", "Use this setup to generate code")
}

type basicModel struct {
	project *project.Project
	key     string
}

type entityModel struct {
	entity  *source.Entity
	project *project.Project
	key     string
}

type entitiesModel struct {
	entity  *[]source.Entity
	project *project.Project
	key     string
}

func (input *entityModel) ToModel() interface{} {
	imports := []string{}

	out := model.EntityModel{}
	if input.entity != nil {
		out = toEntitySection(input.entity)
	}

	if input.project != nil {
		fmt.Println(input.project)
		out.Project.Name = input.project.Name
		out.Project.Owner = input.project.CustomerName

		out.PageHeader = input.project.Header

		if len(input.key) > 0 {
			val, found := input.project.Code.Keys[input.key]
			if found {
				imports = append(imports, val.Imports...)

				out.Inject = make(map[string]model.InjectSection)
				for _, injectKey := range val.Inject {
					injItem, foundInj := input.project.Code.Inject[injectKey]
					if foundInj {
						out.Inject[injectKey] = toInjectSection(injItem)
					}
				}
			}
		}
		if len(input.project.Options) > 0 {
			out.Options = input.project.Options
		}
	}

	return out
}

func (input *basicModel) ToModel() interface{} {
	b := model.BasicModel{}
	imports := []string{}
	// inject := map[]
	if input.project != nil {

		b.Project.Name = input.project.Name
		b.Project.Owner = input.project.CustomerName

		b.PageHeader = input.project.Header

		if len(input.key) > 0 {
			val, found := input.project.Code.Keys[input.key]
			if found {
				imports = append(imports, val.Imports...)

				b.Inject = make(map[string]model.InjectSection)
				for _, injectKey := range val.Inject {
					injItem, foundInj := input.project.Code.Inject[injectKey]
					if foundInj {
						b.Inject[injectKey] = toInjectSection(injItem)
					}
				}
			}
		}
		if len(input.project.Options) > 0 {
			b.Options = input.project.Options
		}
	}

	return b
}

func toColumnSection(from source.Column, entityName string) model.EntityColumnViewModel {
	col := model.EntityColumnViewModel{
		Description:       from.Description,
		IsForeignKey:      from.IsForeignKey,
		IsPrimaryKey:      from.IsPrimaryKey,
		IsIdentity:        from.IsIdentity,
		IsNullable:        from.IsNullable,
		IsIgnored:         from.IsIgnored,
		IsDeletedMarker:   from.IsDeletedMarker,
		IsCreatedDate:     from.IsCreatedDate,
		IsCreatedByUser:   from.IsCreatedByUser,
		IsModifiedDate:    from.IsModifiedByUser,
		IsModifiedByUser:  from.IsModifiedByUser,
		HasPrefix:         strings.HasPrefix(from.Name, entityName),
		HasDescription:    len(from.Description) > 0,
		Name:              from.Name,
		NameWithoutPrefix: strings.TrimPrefix(from.Name, entityName),
		Collation:         from.Collation,
		ReferencesColumn:  from.ReferencesColumn,
		ReferencesTable:   from.ReferencesTable,
		DataType:          from.DataType,
		Length:            from.Length,
		Precision:         from.Precision,
		Scale:             from.Scale,
		UseLength:         from.UseLength,
		UsePrecision:      from.UsePrecision,
	}

	return col
}
func toEntitySection(from *source.Entity) model.EntityModel {
	out := model.EntityModel{
		Name:               from.Name,
		Schema:             from.Schema,
		Type:               from.Type,
		Alias:              from.Alias,
		Description:        from.Description,
		HasDescription:     len(from.Description) > 0,
		HasPrefix:          false,
		NameWithoutPrefix:  "",
		UsesIdentityColumn: from.UsesIdentityColumn,
	}

	for _, column := range from.Columns {

		col := toColumnSection(column, out.Name)

		out.Columns = append(out.Columns, col)

		if column.IsPrimaryKey {
			out.PrimaryKeys = append(out.PrimaryKeys, col)
		}

		if column.IsForeignKey {
			out.ForeignKeys = append(out.ForeignKeys, col)
		}
	}

	for _, cr := range from.ChildRelations {
		child := model.EntityRelationViewModel{}

		child.Name = cr.Name
		child.Schema = cr.Schema
		child.ReleatedColumn = model.EntityColumnProps{
			Name:       cr.ColumnName,
			DataType:   cr.ColumnType,
			IsNullable: cr.ColumnNullable,
		}

		child.OwnerColumn = model.EntityColumnProps{
			Name:       cr.OwnerColumnName,
			DataType:   cr.OwnerColumnType,
			IsNullable: cr.OwnerColumnNullable,
		}

		out.Children = append(out.Children, child)
		// children = append(children)
		out.NameWithoutPrefix = strings.TrimPrefix(cr.Name, out.Name)
		out.HasPrefix = strings.HasPrefix(cr.Name, out.Name)
	}

	for _, pr := range from.ParentRelations {
		parent := model.EntityRelationViewModel{}
		parent.Name = pr.Name
		parent.Schema = pr.Schema

		parent.HasDescription = false
		out.NameWithoutPrefix = strings.TrimPrefix(pr.Name, out.Name)
		out.HasPrefix = strings.HasPrefix(pr.Name, out.Name)
	}

	return out
}
func toInjectSection(from project.CodeInject) model.InjectSection {
	code := model.InjectSection{
		Name:         from.Name,
		PropertyName: from.PropertyName,
	}

	return code
}

// func getCurrentTemplateSet()
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
		Name:              "Contact",
		Description:       "This is a description provided from the table",
		HasDescription:    true,
		HasPrefix:         false,
		NameWithoutPrefix: "Test",
		Columns: []tpl.EntityColumnImportModel{
			{Name: "Id", DataType: "int", IsNullable: false, HasDescription: true, Description: "Description for this column"},
			{Name: "FirstName", DataType: "varchar", IsNullable: false},
			{Name: "LastName", DataType: "varchar", IsNullable: false},
			{Name: "Age", DataType: "int", IsNullable: true},
			{Name: "IsCool", DataType: "bit", IsNullable: true, HasPrefix: true, NameWithoutPrefix: "Cool"},
			{Name: "TypeId", DataType: "int", IsNullable: true},

			{Name: "ChildTest", DataType: "varchar", IsNullable: true, HasPrefix: true, NameWithoutPrefix: "Test", HasDescription: true, Description: "Description for this column"},
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
		Imports:     []string{"using TEST;"},
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
