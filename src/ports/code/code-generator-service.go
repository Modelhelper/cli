package code

import (
	"context"
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	// "go.opencensus.io/examples/exporter"
)

type codeGeneratorService struct {
	// templateService modelhelper.CodeTemplateService
	// app *modelhelper.ModelhelperCli
	cfg             *models.Config
	projectConfig   *models.ProjectConfig
	cmc             modelhelper.CodeModelConverter
	templateService modelhelper.CodeTemplateService
	generator       modelhelper.TemplateGenerator[*models.CodeTemplate]
}

func NewCodeGeneratorService(cfg *models.Config, pc *models.ProjectConfig, cmc modelhelper.CodeModelConverter, ts modelhelper.CodeTemplateService, g modelhelper.TemplateGenerator[*models.CodeTemplate]) modelhelper.CodeGeneratorService {
	return &codeGeneratorService{cfg, pc, cmc, ts, g}
}

func (g *codeGeneratorService) Generate(ctx context.Context, options *models.CodeGeneratorOptions) ([]models.TemplateGeneratorFileResult, error) {
	fmt.Println("Generate code")

	return nil, nil
	// 	if len(options.Templates) == 0 && len(options.TemplateGroups) == 0 {
	// 		// no point to continue if no templates is given

	// 		return nil, errors.New(`No templates or template groups are provided resulting in nothing to create
	// please use mh generate with the -t or --template [templatename] to set at template

	// You could also use mh template or mh t to see a list of all available templates`)
	// 	}

	// 	var con models.Connection
	// 	var prj *models.ProjectConfig
	// 	// var entities []*models.Entity

	// 	if options.UseDemo {
	// 		options.Connection = "demo"
	// 		con = models.Connection{Type: options.Connection}
	// 	} else {
	// 		if len(g.cfg.Connections) == 0 {
	// 			return nil, errors.New("Could not find any connections to use, please add a connection to the config file")
	// 		}
	// 		if len(options.Connection) == 0 {

	// 			options.Connection = g.cfg.DefaultConnection
	// 		}

	// 		if len(options.Connection) == 0 {
	// 			ka := keyArray(g.cfg.Connections)
	// 			options.Connection = ka[0]
	// 		}

	// 		con = g.cfg.Connections[options.Connection]
	// 	}

	// 	entityList := options.Entities //  mergedList(options.Entities, entitiesFromGroups(con, options.EntityGroups))

	// 	src := source.SourceFactory(&con)

	// 	entities, err := src.EntitiesFromNames(entityList)

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	prj = g.projectConfig

	// 	languages, _ := code.LoadFromPath(g.cfg.Languages.Definitions)

	// 	if len(options.TemplatePath) == 0 {
	// 		options.TemplatePath = g.cfg.Templates.Code[0]
	// 	}

	// 	allTemplates, blocks := loadTemplates(options.TemplatePath)

	// 	options.Templates = selectTemplates(allTemplates, options.Templates, options.TemplateGroups)

	// 	start := time.Now()
	// 	var cstat = &models.TemplateGeneratorStatistics{}
	// 	var generatedCode []models.TemplateGeneratorFileResult

	// 	// creates the root context to be passed to each sub routine
	// 	ctxVal := codegen.CodeContextValue{}
	// 	ctxVal.Blocks = blocks

	// 	for _, tname := range options.Templates {

	// 		// var tt *tpl.Template
	// 		// fmt.Println(tname)
	// 		currentTemplate, found := allTemplates[tname]

	// 		if found {
	// 			var codeSection models.Code
	// 			// var csFound = false
	// 			// obsolete when context is completed
	// 			tplMap := make(map[string]string)

	// 			for k, b := range blocks {
	// 				tplMap[k] = b
	// 			}
	// 			tplMap[tname] = currentTemplate.Body

	// 			ctxVal.TemplateName = tname
	// 			ctxVal.Template = currentTemplate.Body

	// 			ctxVal.Datatypes = defaultNoNullDatatype()
	// 			ctxVal.NullableTypes = defaultNullDatatype()

	// 			def, defFound := languages[currentTemplate.Language]
	// 			if defFound {

	// 				for k, v := range def.DataTypes {
	// 					ctxVal.Datatypes[k] = v.NotNull
	// 					ctxVal.NullableTypes[k] = v.Nullable

	// 				}

	// 			}

	// 			// if len(prj.Code. {
	// 			if prj != nil && prj.Code != nil {
	// 				// codeSection, csFound = prj.Code[currentTemplate.Language]
	// 				codeSection = prj.Code[currentTemplate.Language]
	// 			}
	// 			// }
	// 			generator := codegen.GoLangGenerator{}

	// 			ctx := context.WithValue(context.Background(), "code", ctxVal)
	// 			if len(currentTemplate.Model) == 0 || currentTemplate.Model == "basic" {

	// 				basicGenerator := func() {
	// 					cstat.TemplatesUsed += 1

	// 					model := g.cmc.ToBasicModel(currentTemplate.Key, currentTemplate.Language, prj)
	// 					o, _ := generator.Generate(ctx, model)

	// 					f := models.TemplateGeneratorFileResult{
	// 						Result:   o,
	// 						Filename: "",
	// 					}

	// 					generatedCode = append(generatedCode, f)
	// 				}

	// 				basicGenerator()

	// 			} else if currentTemplate.Model == "entity" && len(*entities) > 0 {

	// 				for _, entity := range *entities {

	// 					entityGenerator := func() {
	// 						cstat.TemplatesUsed += 1
	// 						cstat.EntitiesUsed += 1

	// 						model := g.cmc.ToEntityModel(currentTemplate.Key, currentTemplate.Language, prj, &entity)

	// 						model.PageHeader = codegen.Generate("header", model.PageHeader, model)
	// 						model.Namespace = codegen.Generate("namesp", model.Namespace, model)

	// 						for i, imp := range model.Imports {

	// 							model.Imports[i] = codegen.Generate("import", imp, model)
	// 						}

	// 						model.Imports = removeDuplicateStringValues(model.Imports)

	// 						for x, inj := range model.Inject {

	// 							model.Inject[x].Name = codegen.Generate("injprop", inj.Name, model)
	// 						}

	// 						o, _ := generator.Generate(ctx, model)

	// 						// fullPath := ""
	// 						fileName := ""
	// 						if currentTemplate.Type == "file" && len(currentTemplate.FileName) > 0 {
	// 							cstat.FilesCreated += 1

	// 							fileName = codegen.Generate("filename", currentTemplate.FileName, model)

	// 							// if csFound {
	// 							// 	if options.exportPath
	// 							// 	fullPath = filepath.Join(codeSection.Locations[currentTemplate.Key], filen)
	// 							// }
	// 						}

	// 						f := models.TemplateGeneratorFileResult{
	// 							Result:   o,
	// 							Filename: fileName,
	// 							FilePath: codeSection.Locations[currentTemplate.Key],
	// 						}

	// 						generatedCode = append(generatedCode, f)
	// 					}

	// 					entityGenerator()

	// 				}
	// 			} else if currentTemplate.Model == "entities" && len(*entities) > 0 {

	// 				entitiesGenerator := func() {
	// 					cstat.TemplatesUsed += 1
	// 					model := g.cmc.ToEntityListModel(currentTemplate.Key, currentTemplate.Language, prj, entities)
	// 					model.PageHeader = codegen.Generate("header", model.PageHeader, model)

	// 					model.Namespace = codegen.Generate("namesp", model.Namespace, model)

	// 					for i, imp := range model.Imports {

	// 						model.Imports[i] = codegen.Generate("import", imp, model)
	// 					}

	// 					model.Imports = removeDuplicateStringValues(model.Imports)

	// 					for x, inj := range model.Inject {

	// 						model.Inject[x].Name = codegen.Generate("injprop", inj.Name, model)
	// 					}

	// 					o, _ := generator.Generate(ctx, model)

	// 					fileName := ""
	// 					if currentTemplate.Type == "file" && len(currentTemplate.FileName) > 0 {
	// 						cstat.FilesCreated += 1
	// 						fileName = codegen.Generate("filename", currentTemplate.FileName, model)
	// 						// if csFound {

	// 						// 	fullPath = filepath.Join(codeSection.Locations[currentTemplate.Key], filen)
	// 						// }
	// 					}

	// 					f := models.TemplateGeneratorFileResult{
	// 						Result:   o,
	// 						Filename: fileName,
	// 						FilePath: codeSection.Locations[currentTemplate.Key],
	// 					}

	// 					generatedCode = append(generatedCode, f)

	// 				}

	// 				entitiesGenerator()

	// 			}

	// 		}

	// 	}

	// 	sb := strings.Builder{}
	// 	var fwg sync.WaitGroup
	// 	var flock = sync.Mutex{}
	// 	for _, codeBody := range generatedCode {
	// 		cstat.Chars += codeBody.Result.Statistics.Chars
	// 		cstat.Lines += codeBody.Result.Statistics.Lines
	// 		cstat.Words += codeBody.Result.Statistics.Words
	// 		content := []byte(codeBody.Result.Body)
	// 		if options.ExportToScreen {
	// 			screenWriter := exporter.ScreenExporter{}
	// 			screenWriter.Write([]byte(content))
	// 		}

	// 		if options.ExportToClipboard {
	// 			sb.WriteString(string(codeBody.Result.Body))
	// 		}

	// 		if options.ExportByKey && len(codeBody.Filename) > 0 {
	// 			fwg.Add(1)
	// 			go func(filename string, rootPath string, content []byte) {
	// 				defer fwg.Done()
	// 				keyExporter := exporter.FileExporter{
	// 					Filename:  filepath.Join(rootPath, filename),
	// 					Overwrite: options.Overwrite,
	// 				}

	// 				_, err := keyExporter.Write([]byte(content))
	// 				if err != nil {
	// 					fmt.Println(filepath.ErrBadPattern)
	// 				}
	// 				// fmt.Println("*** FILENAME::", s.filename)
	// 				flock.Lock()
	// 				cstat.FilesExported += 1
	// 				flock.Unlock()
	// 			}(codeBody.Filename, "D:/projects/ModelHelper", content)

	// 		}
	// 		// TODO: export to file
	// 		if len(options.ExportPath) > 0 {
	// 			fwg.Add(1)
	// 			go func(filename string, rootPath string, content []byte) {
	// 				defer fwg.Done()

	// 				if len(filename) > 0 {

	// 					fileExporter := exporter.FileExporter{
	// 						Filename:  filepath.Join(rootPath, filename),
	// 						Overwrite: options.Overwrite,
	// 					}

	// 					_, err := fileExporter.Write([]byte(content))
	// 					if err != nil {
	// 						fmt.Printf("%s, err: \n%v", filepath.ErrBadPattern, err)
	// 					}
	// 					// fmt.Println("*** FILENAME::", s.filename)
	// 					flock.Lock()
	// 					cstat.FilesExported += 1
	// 					flock.Unlock()
	// 				} else {
	// 					fmt.Println("Filename empty...")
	// 				}
	// 			}(codeBody.Filename, options.ExportPath, content)
	// 		}
	// 	}

	// 	fwg.Wait()
	// 	if options.ExportToClipboard {
	// 		fmt.Printf("\nGenerated code is copied to the \033[37mclipboard\033[0m. Use \033[34mctrl+v\033[0m to paste it where you like")
	// 		clipboard.WriteAll(sb.String())
	// 	}

	// 	cstat.Duration = time.Since(start)
	// 	// stat["total.time"] = int(cstat.duration.Milliseconds())
	// 	if !options.CodeOnly {
	// 		wpm := 30.0
	// 		cpm := 250.0

	// 		min := float64(cstat.Words) / wpm
	// 		min = float64(cstat.Chars) / cpm
	// 		// stat["total.savings"] = int(min)
	// 		printStat(cstat)
	// 		fmt.Printf("\nIn summary... It took \033[32m%vms\033[0m to generate \033[34m%d\033[0m words and \033[34m%d\033[0m lines. \nYou saved around \033[32m%v minutes\033[0m by not typing it youreself\n",
	// 			cstat.Duration.Milliseconds(),
	// 			cstat.Words,
	// 			cstat.Lines,
	// 			int(min))
	// 	}

	// 	return nil, nil
}

func keyArray(input map[string]models.Connection) []string {
	keys := []string{}
	for k := range input {
		keys = append(keys, k)
	}

	return keys
}

func selectTemplates(templates map[string]models.CodeTemplate, input []string, groups []string) []string {
	list := input

	if len(groups) > 0 {
		for keyTpl, tplVal := range templates {
			for _, templateGroup := range tplVal.Groups {

				for _, checkGroupName := range groups {
					if checkGroupName == templateGroup {
						list = append(list, keyTpl)
					}
				}
			}
		}
	}

	return removeDuplicateStringValues(list)

}

func isInArray(toFind string, items []string) bool {

	for _, entry := range items {
		if entry == toFind {
			return true
		}
	}
	return false
}

func removeDuplicateStringValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
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

// func templatesFromGroups()
func entitiesFromGroups(con models.Connection, groups []string) []string {
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

func printStat(stat *models.TemplateGeneratorStatistics) {
	fmt.Printf(`

Statistics:
---------------------------------------
`)
	tpl := "%-20s%8d\n"

	fmt.Printf(tpl, "Templates used", stat.TemplatesUsed)
	fmt.Printf(tpl, "Entities used", stat.EntitiesUsed)
	fmt.Printf(tpl, "Files created", stat.FilesCreated)
	fmt.Printf(tpl, "Files exported", stat.FilesExported)
	// fmt.Printf(tpl, "Snippets inserted", 1)
	fmt.Println()
	fmt.Printf(tpl, "Character count", stat.Chars)
	fmt.Printf(tpl, "Word count", stat.Words)
	fmt.Printf(tpl, "Line count", stat.Lines)
	fmt.Printf(tpl, "Time used (ms)", stat.Duration.Milliseconds())

}

// func getCurrentTemplateSet()
func testTable() *models.EntityImportModel {
	table := models.EntityImportModel{
		Code: models.CodeImportModel{
			Language: "cs",
			Creator:  models.CreatorImportModel{CompanyName: "FooBar inc", DeveloperName: "Dev E. Loper"},
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
		Columns: []models.EntityColumnImportModel{
			{Name: "Id", DataType: "int", IsNullable: false, HasDescription: true, Description: "Description for this column"},
			{Name: "FirstName", DataType: "varchar", IsNullable: false},
			{Name: "LastName", DataType: "varchar", IsNullable: false},
			{Name: "Age", DataType: "int", IsNullable: true},
			{Name: "IsCool", DataType: "bit", IsNullable: true, HasPrefix: true, NameWithoutPrefix: "Cool"},
			{Name: "TypeId", DataType: "int", IsNullable: true},

			{Name: "ChildTest", DataType: "varchar", IsNullable: true, HasPrefix: true, NameWithoutPrefix: "Test", HasDescription: true, Description: "Description for this column"},
		},
	}
	c1 := models.EntityRelation{
		Name: "ContactAddress",
		ReleatedColumn: models.EntityColumnProps{
			Name: "ContactId", DataType: "int", IsNullable: true,
		},
		OwnerColumn: models.EntityColumnProps{
			Name: "Id", DataType: "int", IsNullable: false,
		},
	}

	p1 := models.EntityRelation{
		Name: "ContactType",
		ReleatedColumn: models.EntityColumnProps{
			Name: "TypeId", DataType: "int", IsNullable: true,
		},
		OwnerColumn: models.EntityColumnProps{
			Name: "Id", DataType: "int", IsNullable: false,
		},
	}

	table.Children = append(table.Children, c1)
	table.Parents = append(table.Parents, p1)
	return &table
}

func testTypes() map[string]models.CodeTypeImportModel {
	tl := make(map[string]models.CodeTypeImportModel)

	tl["model"] = models.CodeTypeImportModel{
		NamePostfix: "",
		NameSpace:   "Testing.Models",
		Key:         "key",
		Imports:     []string{"using HotChocolate;"},
		// Imports:     []string{},
	}
	tl["resolver"] = models.CodeTypeImportModel{
		NamePostfix: "Resolver",
		NameSpace:   "Testing.Resolvers",
		Key:         "key",
		Imports:     []string{"using TEST;"},
	}
	tl["inteface"] = models.CodeTypeImportModel{
		NamePostfix: "Repository",
		NameSpace:   "Testing.Data",
		NamePrefix:  "I",
		Key:         "key",
	}
	tl["repository"] = models.CodeTypeImportModel{
		NamePostfix: "Repository",
		NameSpace:   "Testing.Data",
		Key:         "key",
	}
	return tl
}
