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
	"bufio"
	"fmt"
	"log"
	"modelhelper/cli/app"
	"modelhelper/cli/codegen"
	"modelhelper/cli/config"
	"modelhelper/cli/ctx"
	"modelhelper/cli/model"
	"modelhelper/cli/project"
	"modelhelper/cli/source"
	"modelhelper/cli/tpl"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

type codeContext struct {
	TemplateName             string
	Templates                map[string]string
	Blocks                   map[string]string
	Datatypes                map[string]string
	NullableTypes            map[string]string
	AlternativeNullableTypes map[string]string
}

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
		configFile, _ := cmd.Flags().GetString("config")
		inputTemplates, err := cmd.Flags().GetStringArray("template")
		printScreen, _ := cmd.Flags().GetBool("screen")
		toClipBoard, _ := cmd.Flags().GetBool("copy")

		// if isDemo == false && len(entityFlagArray) == 0 {
		// 	return
		// }
		modelHelperApp = app.New()

		appCtx := modelHelperApp.CreateContext()
		// var ctx *app.Context
		var prj *project.Project
		var entities []source.Entity

		cfg := loadConfig(configFile)

		entities = *loadEntities(*appCtx, entityFlagArray, isDemo)

		prj = loadProject(projectPath)

		if len(tempPath) == 0 {
			tempPath = cfg.Templates.Location
		}

		allTemplates, blocks := loadTemplates(tempPath)

		if err != nil {
			panic(err)
		}

		start := time.Now()
		var cstat = codeStat{}
		var generatedCode []codeFile
		if len(inputTemplates) > 0 {

			for _, tname := range inputTemplates {

				// var tt *tpl.Template
				// fmt.Println(tname)
				currentTemplate, found := allTemplates[tname]

				if found {

					var input interface{}

					// obsolete when context is completed
					tplMap := make(map[string]string)

					for k, b := range blocks {
						tplMap[k] = b
					}
					tplMap[tname] = currentTemplate.Body

					// create context
					// obsolete when context is completed
					codeCtx := ctx.Context{}
					codeCtx.TemplateName = tname
					codeCtx.Templates = blocks

					// obsolete when context is completed
					generator := codegen.GoLangGenerator{
						Templates:    tplMap,
						TemplateName: tname,
					}

					// TODO: refactor to method
					if len(currentTemplate.Model) == 0 || currentTemplate.Model == "basic" {
						tstart := time.Now()

						basicModel := basicModel{
							project: prj,
							key:     currentTemplate.Key,
						}

						input = basicModel.ToModel(codeCtx)
						o, _ := generator.Generate(codeCtx, input)
						f := codeFile{
							content: o,
							stat:    getStat(o),
						}
						f.stat.duration = time.Since(tstart)

						generatedCode = append(generatedCode, f)
					}

					// TODO: refactor to method
					if currentTemplate.Model == "entity" && len(entities) > 0 {

						for _, entity := range entities {
							// stat["entities"] += 1
							tstart := time.Now()

							entityModel := entityModel{
								entity:  &entity,
								project: prj,
								key:     currentTemplate.Key,
							}
							input = entityModel.ToModel(codeCtx)

							// input = getEntityModel(codeCtx, &entity, prj, currentTemplate.Key)
							// f := generateCode(codeCtx, input)
							o, _ := generator.Generate(codeCtx, input)

							f := codeFile{
								content: o,
								stat:    getStat(o),
							}

							f.stat.duration = time.Since(tstart)

							generatedCode = append(generatedCode, f)

						}
					}

				}

			}

			sb := strings.Builder{}
			for _, s := range generatedCode {
				cstat.appendStat(s.stat)

				if printScreen {
					screenWriter := tpl.ScreenExporter{}
					screenWriter.Export([]byte(s.content))
				}

				if toClipBoard {
					sb.WriteString(s.content)
				}
				// TODO: export to file
			}

			if toClipBoard {
				fmt.Printf("\nGenerated code is copied to the \033[37mclipboard\033[0m. Use \033[34mctrl+v\033[0m to paste it where you like")
				clipboard.WriteAll(sb.String())
			}

			cstat.duration = time.Since(start)
			// stat["total.time"] = int(cstat.duration.Milliseconds())
			if !codeOnly {
				wpm := 40.0
				min := float64(cstat.words) / wpm
				// stat["total.savings"] = int(min)
				printStat(cstat)
				fmt.Printf("\nIn summary... It took \033[32m%vms\033[0m to generate \033[34m%d\033[0m words and \033[34m%d\033[0m lines. \nYou saved around \033[32m%v minutes\033[0m by not typing it youreself\n",
					cstat.duration.Milliseconds(),
					cstat.words,
					cstat.lines,
					int(min))
			}
		}

	},
}

func (cstat *codeStat) appendStat(instat codeStat) {
	cstat.chars += instat.chars
	cstat.lines += instat.lines
	cstat.words += instat.words
}
func getEntityModel(c ctx.Context, entity *source.Entity, prj *project.Project, key string) interface{} {
	mdl := entityModel{
		entity:  entity,
		project: prj,
		key:     key,
	}
	input := mdl.ToModel(c)

	return input
}
func generateCode(cctx ctx.Context, model interface{}) codeFile {
	generator := codegen.GoLangGenerator{}
	start := time.Now()

	o, _ := generator.Generate(cctx, model)
	f := codeFile{
		content: o,
		stat:    getStat(o),
	}
	f.stat.duration = time.Since(start)

	return f
}

func getStat(s string) codeStat {
	stat := codeStat{
		chars: len(s),
		lines: getLines(s),
		words: getWords(s),
	}

	return stat
}

type codeFile struct {
	filename        string
	content         string
	stat            codeStat
	exists          bool
	existingContent string
}

type codeStat struct {
	chars     int
	lines     int
	words     int
	duration  time.Duration
	timeSaved int
}

func getWords(input string) int {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}

	return count
}
func getLines(input string) int {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanLines)

	count := 0
	for scanner.Scan() {
		count++
	}

	return count
}
func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringArrayP("template", "t", []string{}, "a list of template to convert")
	generateCmd.Flags().StringArrayP("entity", "e", []string{}, "a list of entits to use as a model")
	generateCmd.Flags().Bool("screen", false, "List the output to the screen, default false")
	generateCmd.Flags().Bool("copy", false, "Copies the generated code to the clipboard (ctrl + v), default false")
	generateCmd.Flags().String("export", "", "Exports to a directory")
	generateCmd.Flags().Bool("export-bykey", false, "Exports the code using the template keys, default false")
	generateCmd.Flags().Bool("code-only", false, "Writes only the generated code to the console, no stats, no messages - only code, default false")
	generateCmd.Flags().Bool("demo", false, "Uses a demo as input source, this will override any other input sources (entity, graphql), default false ")

	generateCmd.Flags().String("template-path", "", "Instructs the program to use this path as root for templates")
	generateCmd.Flags().String("config", "", "Instructs the program to use this config as the config")
	generateCmd.Flags().String("project", "", "Instructs the program to use this project as input")

	generateCmd.Flags().String("key", "", "The key to use when encoding and decoding secrets for a connection")

	generateCmd.Flags().String("setup", "", "Use this setup to generate code")
}

func printStat(stat codeStat) {
	fmt.Printf(`

Statistics:
---------------------------------------
`)
	tpl := "%-20s%8d\n"

	fmt.Printf(tpl, "Templates used", 2)
	fmt.Printf(tpl, "Entities used", 4)
	fmt.Printf(tpl, "Files exported", 6)
	fmt.Printf(tpl, "Snippets inserted", 1)
	fmt.Println()
	fmt.Printf(tpl, "Character count", stat.chars)
	fmt.Printf(tpl, "Word count", stat.words)
	fmt.Printf(tpl, "Line count", stat.lines)
	fmt.Printf(tpl, "Time used (ms)", stat.duration.Milliseconds())

}

type basicModel struct {
	project   *project.Project
	generator *codegen.SimpleGenerator
	key       string
}

type entityModel struct {
	entity    *source.Entity
	project   *project.Project
	generator *codegen.SimpleGenerator
	key       string
}

type entitiesModel struct {
	entity    *[]source.Entity
	project   *project.Project
	generator *codegen.SimpleGenerator
	key       string
}

func loadTemplates(templatePath string) (all map[string]tpl.Template, blocks map[string]string) {

	tl := tpl.TemplateLoader{
		Directory: app.TemplateFolder(templatePath),
	}

	a, _ := tl.LoadTemplates()

	var b = make(map[string]string)

	for tk, tv := range a {
		if strings.ToLower(tv.Type) == "block" {
			b[tk] = tv.Body
		}
	}
	// b := tpl.ExtractBlocks(&a)

	return a, b
}

func loadEntities(appCtx app.Context, names []string, isDemo bool) *[]source.Entity {
	var entities []source.Entity

	if isDemo {
		// load demo project
		var ds *source.DemoSource
		el, _ := ds.Entities("")

		for _, eitem := range *el {
			entities = append(entities, eitem)
		}

		// load demo tables (2)
	} else {

		conName := appCtx.DefaultConnection

		con := appCtx.Connections[conName]
		src := con.LoadSource()

		if len(names) > 0 {
			for _, entityName := range names {
				entity, err := src.Entity(entityName)
				if err != nil {
					log.Fatalln(err)
				}

				entities = append(entities, *entity)
			}
		}

	}

	return &entities
}

func loadConfig(configPath string) *config.Config {
	if len(configPath) > 0 {
		return config.LoadFromFile(configPath)
	} else {
		return config.Load()
	}
}
func loadProject(projectPath string) *project.Project {
	if len(projectPath) == 0 {
		if project.Exists(project.DefaultLocation()) {
			projectPath = project.DefaultLocation()
		} else {
			fp, foundProject := project.FindNearestProjectDir()
			if foundProject && project.Exists(fp) {
				projectPath = fp

			}
		}
	}

	if len(projectPath) > 0 {

		fmt.Println("Project found here: ", projectPath)
		prj, _ := project.Load(projectPath)
		return prj
	}

	return nil
}
func (input *entityModel) ToModel(codeCtx ctx.Context) interface{} {

	imports := []string{}

	out := model.EntityModel{}

	if input.entity != nil {
		out = toEntitySection(input.entity)
	}

	if input.project != nil {
		if len(input.project.Options) > 0 {
			out.Options = input.project.Options
		}
		out.Project.Name = input.project.Name
		out.Project.Owner = input.project.OwnerName

		out.PageHeader = input.project.Header
		if len(input.key) > 0 {
			val, found := input.project.Code.Keys[input.key]
			if found {
				out.Namespace = val.NameSpace
				out.Postfix = val.Postfix
				out.Prefix = val.Prefix

				for _, imp := range val.Imports {
					imports = append(imports, imp)
				}
				out.Imports = imports
				out.Inject = []model.InjectSection{}
				for _, injectKey := range val.Inject {
					injItem, foundInj := input.project.Code.Inject[injectKey]
					if foundInj {
						for _, injImport := range injItem.Imports {
							out.Imports = append(out.Imports, injImport)
						}
						out.Inject = append(out.Inject, toInjectSection(injItem, out))

					}
				}
			}
		}

	}

	return out
}

func (input *basicModel) ToModel(codeCtx ctx.Context) interface{} {
	b := model.BasicModel{}
	imports := []string{}
	// inject := map[]
	if input.project != nil {

		if len(input.project.Options) > 0 {
			fmt.Println("has options")
			b.Options = input.project.Options
		}

		b.Project.Name = input.project.Name
		b.Project.Owner = input.project.OwnerName

		b.PageHeader = input.project.Header

		if len(input.key) > 0 {
			val, found := input.project.Code.Keys[input.key]
			if found {
				imports = append(imports, val.Imports...)

				b.Inject = []model.InjectSection{}
				for _, injectKey := range val.Inject {
					injItem, foundInj := input.project.Code.Inject[injectKey]
					if foundInj {
						b.Inject = append(b.Inject, toInjectSection(injItem, b))
					}
				}
			}
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
func toInjectSection(from project.CodeInject, m interface{}) model.InjectSection {
	name, _ := codegen.Generate("fileName", from.Name, m)
	code := model.InjectSection{
		Name:         name,
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
