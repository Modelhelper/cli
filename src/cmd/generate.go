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
	"context"
	"fmt"
	"log"
	"modelhelper/cli/app"
	"modelhelper/cli/code"
	"modelhelper/cli/codegen"
	"modelhelper/cli/config"
	"modelhelper/cli/model"
	"modelhelper/cli/project"
	"modelhelper/cli/source"
	"modelhelper/cli/tpl"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringArrayP("template", "t", []string{}, "A list of template to convert")
	generateCmd.Flags().StringArray("template-group", []string{}, "Use a group of templates")
	generateCmd.Flags().String("template-path", "", "Instructs the program to use this path as root for templates")

	generateCmd.Flags().StringP("relations [direct, all, complete]", "r", "", "Include related entities based on the entities in --entity or --entity-group ('direct' | 'all' | 'complete' | 'children' | 'parents')")
	// generateCmd.Flags().String("template-path", "", "Instructs the program to use this path as root for templates")

	generateCmd.Flags().StringArray("entity-group", []string{}, "Use a group of entities (must be defines in the current connection)")
	generateCmd.Flags().StringArrayP("entity", "e", []string{}, "A list of entits to use as a model")

	generateCmd.Flags().Bool("screen", false, "List the output to the screen, default false")
	generateCmd.Flags().Bool("copy", false, "Copies the generated code to the clipboard (ctrl + v), default false")
	generateCmd.Flags().String("export-path", "", "Exports to a directory")
	generateCmd.Flags().Bool("export-bykey", false, "Exports the code using the template keys, default false")
	generateCmd.Flags().Bool("overwrite", false, "Overwrite any existing file when exporting to file on disk")

	generateCmd.Flags().Bool("code-only", false, "Writes only the generated code to the console, no stats, no messages - only code, default false")

	generateCmd.Flags().Bool("demo", false, "Uses a demo as input source, this will override any other input sources (entity, graphql), default false ")

	generateCmd.Flags().String("config-path", "", "Instructs the program to use this config as the config")
	generateCmd.Flags().String("project-path", "", "Instructs the program to use this project as input")

	generateCmd.Flags().String("key", "", "The key to use when encoding and decoding secrets for a connection")

	// generateCmd.Flags().String("setup", "", "Use this setup to generate code") // version 3.1
	generateCmd.Flags().StringP("connection", "c", "", "The connection key to be used, uses default connection if not provided")

	generateCmd.RegisterFlagCompletionFunc("relations", completeRelations)

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
		entityGroupFlagArray, _ := cmd.Flags().GetStringArray("entity-group")
		tempPath, _ := cmd.Flags().GetString("template-path")
		projectPath, _ := cmd.Flags().GetString("project-path")
		configFile, _ := cmd.Flags().GetString("config-path")
		inputTemplates, err := cmd.Flags().GetStringArray("template")
		inputGroupTemplates, err := cmd.Flags().GetStringArray("template-group")
		printScreen, _ := cmd.Flags().GetBool("screen")
		toClipBoard, _ := cmd.Flags().GetBool("copy")
		exportByKey, _ := cmd.Flags().GetBool("export-bykey")
		conName, _ := cmd.Flags().GetString("connection")
		overwriteAll, _ := cmd.Flags().GetBool("overwrite")

		if len(inputTemplates) == 0 {
			// no point to continue if no templates is given
			fmt.Printf(`No templates or template groups are provided resulting in nothing to create
please use mh generate with the -t or --template [templatename] to set at template

You could also use mh template or mh t to see a list of all available templates`)

			return
		}

		var wg sync.WaitGroup
		var lock = sync.Mutex{}

		// obsolete
		modelHelperApp = app.New()

		// obsolete
		appCtx := modelHelperApp.CreateContext()

		if len(conName) == 0 {

			conName = appCtx.DefaultConnection
		}

		if len(conName) == 0 {
			ka := keyArray(appCtx.Connections)
			conName = ka[0]
		}

		var prj *project.Project
		var entities []source.Entity

		cfg := loadConfig(configFile)

		con := appCtx.Connections[conName]

		entityList := mergedList(entityFlagArray, entitiesFromGroups(con, entityGroupFlagArray))

		src := con.LoadSource()

		entities = *loadEntities(src, entityList, isDemo)

		prj = loadProject(projectPath)

		languages, _ := code.LoadFromPath(cfg.Languages.Definitions)

		if len(tempPath) == 0 {
			tempPath = cfg.Templates.Location
		}

		allTemplates, blocks := loadTemplates(tempPath)

		if err != nil {
			panic(err)
		}

		start := time.Now()
		var cstat = codegen.Statistics{}
		var generatedCode []codeFile

		// creates the root context to be passed to each sub routine
		ctxVal := codegen.CodeContextValue{}
		ctxVal.Blocks = blocks

		for _, tname := range inputTemplates {

			// var tt *tpl.Template
			// fmt.Println(tname)
			currentTemplate, found := allTemplates[tname]

			if found {
				var codeSection code.Code
				var csFound = false
				// obsolete when context is completed
				tplMap := make(map[string]string)

				for k, b := range blocks {
					tplMap[k] = b
				}
				tplMap[tname] = currentTemplate.Body

				ctxVal.TemplateName = tname
				ctxVal.Template = currentTemplate.Body

				ctxVal.Datatypes = defaultNoNullDatatype()
				ctxVal.NullableTypes = defaultNullDatatype()

				def, defFound := languages[currentTemplate.Language]
				if defFound {

					for k, v := range def.DataTypes {
						ctxVal.Datatypes[k] = v.NotNull
						ctxVal.NullableTypes[k] = v.Nullable

					}

				}

				// if len(prj.Code. {
				if prj != nil && prj.Code != nil {
					codeSection, csFound = prj.Code[currentTemplate.Language]
				}
				// }
				generator := codegen.GoLangGenerator{}

				ctx := context.WithValue(context.Background(), "code", ctxVal)
				if len(currentTemplate.Model) == 0 || currentTemplate.Model == "basic" {

					basicGenerator := func() {
						defer wg.Done()

						model := ToBasicModel(currentTemplate.Key, currentTemplate.Language, prj)
						o, _ := generator.Generate(ctx, model)

						f := codeFile{
							result:   o,
							filename: "",
						}

						lock.Lock()
						generatedCode = append(generatedCode, f)
						lock.Unlock()
					}

					wg.Add(1)
					go basicGenerator()

				} else if currentTemplate.Model == "entity" && len(entities) > 0 {

					for _, entity := range entities {

						entityGenerator := func() {
							defer wg.Done()
							model := ToEntityModel(currentTemplate.Key, currentTemplate.Language, prj, &entity)

							model.PageHeader = codegen.Generate("header", model.PageHeader, model)
							model.Namespace = codegen.Generate("namesp", model.Namespace, model)

							for i, imp := range model.Imports {

								model.Imports[i] = codegen.Generate("import", imp, model)
							}

							for x, inj := range model.Inject {

								model.Inject[x].Name = codegen.Generate("injprop", inj.Name, model)
							}

							o, _ := generator.Generate(ctx, model)
							filen := codegen.Generate("filename", currentTemplate.FileName, model)

							fullPath := ""
							if csFound {

								fullPath = filepath.Join(codeSection.Locations[currentTemplate.Key], filen)
							}

							f := codeFile{
								result:   o,
								filename: fullPath,
							}

							lock.Lock()
							generatedCode = append(generatedCode, f)
							lock.Unlock()
						}

						wg.Add(1)
						go entityGenerator()

					}
				} else if currentTemplate.Model == "entities" && len(entities) > 0 {

					entitiesGenerator := func() {
						defer wg.Done()
						model := ToEntitiesModel(currentTemplate.Key, currentTemplate.Language, prj, &entities)
						model.PageHeader = codegen.Generate("header", model.PageHeader, model)

						o, _ := generator.Generate(ctx, model)
						filen := codegen.Generate("filename", currentTemplate.FileName, model)
						fullPath := ""
						if csFound {

							fullPath = filepath.Join(codeSection.Locations[currentTemplate.Key], filen)
						}

						f := codeFile{
							result:   o,
							filename: fullPath,
						}

						lock.Lock()
						generatedCode = append(generatedCode, f)
						lock.Unlock()

					}

					wg.Add(1)
					go entitiesGenerator()

				}

			}

		}

		wg.Wait()

		sb := strings.Builder{}
		for _, s := range generatedCode {
			cstat.AppendStat(s.result.Stat)

			if printScreen {
				screenWriter := tpl.ScreenExporter{}
				screenWriter.Write([]byte(s.result.Content))
			}

			if toClipBoard {
				sb.WriteString(s.result.Content)
			}

			fmt.Println("*** FILENAME::", s.filename)
			// TODO: export to file
		}

		if toClipBoard {
			fmt.Printf("\nGenerated code is copied to the \033[37mclipboard\033[0m. Use \033[34mctrl+v\033[0m to paste it where you like")
			clipboard.WriteAll(sb.String())
		}

		cstat.Duration = time.Since(start)
		// stat["total.time"] = int(cstat.duration.Milliseconds())
		if !codeOnly {
			wpm := 40.0
			min := float64(cstat.Words) / wpm
			// stat["total.savings"] = int(min)
			printStat(cstat)
			fmt.Printf("\nIn summary... It took \033[32m%vms\033[0m to generate \033[34m%d\033[0m words and \033[34m%d\033[0m lines. \nYou saved around \033[32m%v minutes\033[0m by not typing it youreself\n",
				cstat.Duration.Milliseconds(),
				cstat.Words,
				cstat.Lines,
				int(min))
		}

	},
}

func defaultNoNullDatatype() map[string]string {
	// build current template context
	dtm := make(map[string]string)
	dtm["varchar"] = "string"
	dtm["nvarchar"] = "string"
	dtm["datetimeoffset"] = "DateTimeOffset"
	dtm["datetime2"] = "DateTimeOffset"
	dtm["bit"] = "bool"
	dtm["decimal"] = "decimal"

	return dtm

}

func defaultNullDatatype() map[string]string {

	ndtm := make(map[string]string)
	ndtm["varchar"] = "string"
	ndtm["nvarchar"] = "string"
	ndtm["int"] = "int?"
	ndtm["datetimeoffset"] = "DateTimeOffset?"
	ndtm["datetime2"] = "DateTimeOffset?"
	ndtm["bit"] = "bool?"
	ndtm["decimal"] = "decimal?"

	return ndtm
}

func entitiesFromGroups(con source.Connection, groups []string) []string {
	list := []string{}

	for _, group := range groups {

		conGrp, found := con.Groups[group]
		if found {
			for _, e := range conGrp.Items {
				list = append(list, e)
			}
		}
	}

	return list
}

func mergedList(lists ...[]string) []string {
	items := make(map[string]int)
	out := []string{}

	for _, list := range lists {
		for _, item := range list {
			items[item] += 1
		}
	}

	for key, _ := range items {
		out = append(out, key)
	}

	return out

}

type codeFile struct {
	filename        string
	result          codegen.Result
	exists          bool
	existingContent string
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringArrayP("template", "t", []string{}, "A list of template to convert")
	generateCmd.Flags().StringArray("template-group", []string{}, "Use a group of templates")
	generateCmd.Flags().String("template-path", "", "Instructs the program to use this path as root for templates")

	generateCmd.Flags().StringP("relations [direct, all, complete]", "r", "", "Include related entities based on the entities in --entity or --entity-group ('direct' | 'all' | 'complete' | 'children' | 'parents')")
	// generateCmd.Flags().String("template-path", "", "Instructs the program to use this path as root for templates")

	generateCmd.Flags().StringArray("entity-group", []string{}, "Use a group of entities (must be defines in the current connection)")
	generateCmd.Flags().StringArrayP("entity", "e", []string{}, "A list of entits to use as a model")

	generateCmd.Flags().Bool("screen", false, "List the output to the screen, default false")
	generateCmd.Flags().Bool("copy", false, "Copies the generated code to the clipboard (ctrl + v), default false")
	generateCmd.Flags().String("export-path", "", "Exports to a directory")
	generateCmd.Flags().Bool("export-bykey", false, "Exports the code using the template keys, default false")

	generateCmd.Flags().Bool("code-only", false, "Writes only the generated code to the console, no stats, no messages - only code, default false")

	generateCmd.Flags().Bool("demo", false, "Uses a demo as input source, this will override any other input sources (entity, graphql), default false ")

	generateCmd.Flags().String("config-path", "", "Instructs the program to use this config as the config")
	generateCmd.Flags().String("project-path", "", "Instructs the program to use this project as input")

	generateCmd.Flags().String("key", "", "The key to use when encoding and decoding secrets for a connection")

	// generateCmd.Flags().String("setup", "", "Use this setup to generate code") // version 3.1
	generateCmd.Flags().StringP("connection", "c", "", "The connection key to be used, uses default connection if not provided")

	generateCmd.RegisterFlagCompletionFunc("relations", completeRelations)

}

func completeRelations(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"direct", "all", "complete", "children", "parents"}, cobra.ShellCompDirectiveDefault
}
func printStat(stat codegen.Statistics) {
	fmt.Printf(`

Statistics:
---------------------------------------
`)
	tpl := "%-20s%8d\n"

	// fmt.Printf(tpl, "Templates used", 2)
	// fmt.Printf(tpl, "Entities used", 4)
	// fmt.Printf(tpl, "Files exported", 6)
	// fmt.Printf(tpl, "Snippets inserted", 1)
	fmt.Println()
	fmt.Printf(tpl, "Character count", stat.Chars)
	fmt.Printf(tpl, "Word count", stat.Words)
	fmt.Printf(tpl, "Line count", stat.Lines)
	fmt.Printf(tpl, "Time used (ms)", stat.Duration.Milliseconds())

}

type basicModel struct {
	project  *project.Project
	key      string
	language string
}

type entityModel struct {
	entity   *source.Entity
	project  *project.Project
	key      string
	language string
}

type entitiesModel struct {
	entities *[]source.Entity
	project  *project.Project
	key      string
	// generator *codegen.SimpleGenerator
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

func loadEntities(src source.Source, names []string, isDemo bool) *[]source.Entity {
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

		// conName := appCtx.DefaultConnection

		// con := appCtx.Connections[conName]
		// src := con.LoadSource()

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

		prj, _ := project.Load(projectPath)
		return prj
	}

	return nil
}

func toProjectModel(input project.Project) model.ProjectSection {
	return model.ProjectSection{
		Owner: input.OwnerName,
		Name:  input.Name,
	}

}

func ToEntityModel(key, language string, project *project.Project, entity *source.Entity) model.EntityModel {

	entityBase := model.EntityModel{}

	base := ToBasicModel(key, language, project)
	if entity != nil {
		entityBase = toEntitySection(entity)
	}

	out := model.EntityModel{
		RootNamespace:             base.RootNamespace,
		Namespace:                 base.Namespace,
		Postfix:                   base.Postfix,
		Prefix:                    base.Prefix,
		ModuleLevelVariablePrefix: base.ModuleLevelVariablePrefix,
		Inject:                    base.Inject,
		Imports:                   base.Imports,
		Project:                   base.Project,
		Developer:                 base.Developer,
		Options:                   base.Options,
		PageHeader:                base.PageHeader,
		Name:                      entityBase.Name,
		Schema:                    entityBase.Schema,
		Type:                      entityBase.Type,
		Alias:                     entityBase.Alias,
		Description:               entityBase.Description,
		HasDescription:            len(entityBase.Description) > 0,
		HasPrefix:                 false, //len(entityBase.Prefix) > 0,
		NameWithoutPrefix:         "",
		Columns:                   entityBase.Columns,
		Parents:                   entityBase.Parents,
		Children:                  entityBase.Children,
		PrimaryKeys:               entityBase.PrimaryKeys,
		ForeignKeys:               entityBase.ForeignKeys,
		UsedAsColumns:             entityBase.UsedAsColumns,
		UsesIdentityColumn:        entityBase.UsesIdentityColumn,
	}

	return out
}
func ToEntitiesModel(key, language string, project *project.Project, entities *[]source.Entity) model.EntityListModel {

	entitylist := []model.EntityModel{}

	base := ToBasicModel(key, language, project)
	if entities != nil && len(*entities) > 0 {
		for _, entity := range *entities {

			entityBase := toEntitySection(&entity)
			entitylist = append(entitylist, entityBase)
		}
	}

	out := model.EntityListModel{
		Namespace:                 base.Namespace,
		Postfix:                   base.Postfix,
		Prefix:                    base.Prefix,
		ModuleLevelVariablePrefix: base.ModuleLevelVariablePrefix,
		Inject:                    base.Inject,
		Imports:                   base.Imports,
		Project:                   base.Project,
		Developer:                 base.Developer,
		Options:                   base.Options,
		PageHeader:                base.PageHeader,
		Entities:                  entitylist,
	}

	return out
}

func ToBasicModel(key, language string, project *project.Project) model.BasicModel {
	b := model.BasicModel{}
	imports := []string{}

	// inject := map[]
	if project != nil {
		code, codeFound := project.Code[language]

		if len(project.Options) > 0 {
			b.Options = project.Options
		}
		b.Project.Name = project.Name
		b.Project.Owner = project.OwnerName

		b.PageHeader = project.Header

		if len(key) > 0 && codeFound {
			val, found := code.Keys[key]
			if found {
				b.RootNamespace = code.RootNamespace
				imports = append(imports, val.Imports...)

				b.Inject = []model.InjectSection{}
				for _, injectKey := range val.Inject {
					injItem, foundInj := code.Inject[injectKey]
					if foundInj {
						b.Inject = append(b.Inject, toInjectSection(injItem, b))
					}
				}

				b.Postfix = val.Postfix
				b.Prefix = val.Prefix
				b.Namespace = val.Namespace

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
func toInjectSection(from code.Inject, m interface{}) model.InjectSection {
	// name, _ := codegen.Generate("fileName", from.Name, m)
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
