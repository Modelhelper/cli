package code

import (
	"context"
	"errors"
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"os"
	"path/filepath"
	"time"

	"encoding/json"
	// "go.opencensus.io/examples/exporter"
)

type codeGeneratorService struct {
	// templateService modelhelper.CodeTemplateService
	// app *modelhelper.ModelhelperCli
	cfg               *models.Config
	projectConfig     *models.ProjectConfig
	cmc               modelhelper.CodeModelConverter
	templateService   modelhelper.CodeTemplateService
	generator         modelhelper.TemplateGenerator[*models.CodeTemplate]
	connectionService modelhelper.ConnectionService
	sourceFactory     modelhelper.SourceFactoryService
	commitHistory     modelhelper.CommitHistoryService
}

func NewCodeGeneratorService(cfg *models.Config,
	pc *models.ProjectConfig, cmc modelhelper.CodeModelConverter,
	ts modelhelper.CodeTemplateService, g modelhelper.TemplateGenerator[*models.CodeTemplate],
	c modelhelper.ConnectionService, srcf modelhelper.SourceFactoryService,
	ch modelhelper.CommitHistoryService,
) modelhelper.CodeGeneratorService {
	return &codeGeneratorService{cfg, pc, cmc, ts, g, c, srcf, ch}
}

func (g *codeGeneratorService) Generate(ctx context.Context, options *models.CodeGeneratorOptions) (*models.CodeGenerateResult, error) {

	result := new(models.CodeGenerateResult)
	if len(options.Templates) == 0 && len(options.FeatureTemplates) == 0 {
		// no point to continue if no templates is given

		return nil, errors.New(`no templates or template groups are provided resulting in nothing to create
please use mh generate with the -t or --template [templatename] to set at template

You could also use mh template or mh t to see a list of all available templates`)
	}

	// var con models.Connection
	var prj *models.ProjectConfig
	// var entities []*models.Entity
	conType, conName := "", ""

	if options.UseDemo {
		options.ConnectionName = "demo"
		conName = "demo"
		conType = "file"
		// con = models.Connection{Type: options.ConnectionName}
		// con = models.Connection{Type: options.ConnectionName}
	} else {

		connections, err := g.connectionService.Connections()
		if err != nil {
			return nil, err
		}
		if len(connections) == 0 {
			return nil, errors.New("could not find any connections to use, please add a connection to the config file")
		}
		if len(options.ConnectionName) == 0 {

			options.ConnectionName = g.cfg.DefaultConnection
		}

		if len(options.ConnectionName) == 0 {
			for _, v := range connections {
				options.ConnectionName = v.Name
				break
			}
		}

		conName = options.ConnectionName
		conType = connections[conName].Type
	}

	con, _ := g.connectionService.BaseConnection(options.ConnectionName)

	// entityList := options.Entities //  mergedList(options.Entities, entitiesFromGroups(con, options.EntityGroups))
	entityList := mergedList(options.SourceItems, entitiesFromGroups(con, options.SourceItemGroups))

	src, _ := g.sourceFactory.NewSource(conType, conName)

	entities, err := src.EntitiesFromNames(entityList)

	if err != nil {
		return nil, err
	}

	prj = g.projectConfig

	// 	if len(options.TemplatePath) == 0 {
	// 		options.TemplatePath = g.cfg.Templates.Code[0]
	// 	}
	templateOptions := &models.CodeTemplateListOptions{
		DatabaseType: con.Type,
	}

	allTemplates := g.templateService.List(templateOptions)

	options.Templates = selectTemplates(allTemplates, options.Templates, options.FeatureTemplates)

	if len(options.Templates) > 0 && options.Verbose {
		fmt.Println("Templates to use:")
		for _, t := range options.Templates {
			fmt.Println("\t", t)
		}
	}

	start := time.Now()
	var cstat = &models.TemplateGeneratorStatistics{}
	var generatedCode []models.TemplateGeneratorFileResult
	var ch *models.CommitHistory

	for _, tname := range options.Templates {

		// var tt *tpl.Template
		// fmt.Println(tname)
		currentTemplate, found := allTemplates[tname]

		if found {
			if options.Verbose {
				fmt.Printf("Generating template '%s' for '%s' \n", currentTemplate.Name, currentTemplate.Language)
			}
			locationPath, locationFound := "", false

			if prj != nil && prj.Locations != nil {
				locationPath, locationFound = prj.Locations[currentTemplate.Key]
			}

			// var codeSection models.Code

			// if prj != nil && prj.Code != nil {
			// 	codeSection = prj.Code[currentTemplate.Language]
			// }

			if len(currentTemplate.Model) == 0 || currentTemplate.Model == "basic" {

				basicGenerator := func() {
					cstat.TemplatesUsed += 1

					model := g.cmc.ToBasicModel(currentTemplate.Key, currentTemplate.Language, prj)
					o, _ := g.generator.Generate(ctx, &currentTemplate, model, templateOptions)
					code := createCodeResultFile(o, &currentTemplate, model, locationFound, locationPath)
					appendCodeResult(result, code)

				}

				basicGenerator()

			} else if currentTemplate.Model == "name" { //&& len(*entities) > 0
				nameGenerator := func() {
					cstat.TemplatesUsed += 1

					model := g.cmc.ToNameModel(currentTemplate.Key, currentTemplate.Language, prj, options.Name)
					// model.PageHeader = simpleGenerate("header", model.PageHeader, model)

					o, _ := g.generator.Generate(ctx, &currentTemplate, model, templateOptions)
					code := createCodeResultFile(o, &currentTemplate, model, locationFound, locationPath)
					appendCodeResult(result, code)

				}

				nameGenerator()
			} else if currentTemplate.Model == "entity" { //&& len(*entities) > 0

				for _, entity := range *entities {

					if options.Verbose {
						fmt.Printf("\tUsing entity '%s' \n", entity.Name)
					}

					entityGenerator := func() {
						cstat.TemplatesUsed += 1
						cstat.EntitiesUsed += 1

						model := g.cmc.ToEntityModel(currentTemplate.Key, currentTemplate.Language, prj, &entity)

						if options.Verbose {
							fmt.Printf("\tEntity model created: '%s' \n", model.Name)
						}

						model.PageHeader = simpleGenerate("header", model.PageHeader, model)
						model.Namespace = simpleGenerate("namesp", model.Namespace, model)

						for i, imp := range model.Imports {

							model.Imports[i] = simpleGenerate("import", imp, model)
						}

						model.Imports = removeDuplicateStringValues(model.Imports)

						for x, inj := range model.Inject {

							model.Inject[x].Name = simpleGenerate("injprop", inj.Name, model)
						}

						o, err := g.generator.Generate(ctx, &currentTemplate, model, templateOptions)

						if err != nil {
							fmt.Println("Error when generating", err)
						}

						code := createCodeResultFile(o, &currentTemplate, model, locationFound, locationPath)
						appendCodeResult(result, code)

					}

					entityGenerator()

				}
			} else if currentTemplate.Model == "entities" && len(*entities) > 0 {

				entitiesGenerator := func() {
					cstat.TemplatesUsed += 1
					model := g.cmc.ToEntityListModel(currentTemplate.Key, currentTemplate.Language, prj, entities)
					model.PageHeader = simpleGenerate("header", model.PageHeader, model)

					model.Namespace = simpleGenerate("namesp", model.Namespace, model)

					for i, imp := range model.Imports {

						model.Imports[i] = simpleGenerate("import", imp, model)
					}

					model.Imports = removeDuplicateStringValues(model.Imports)

					for x, inj := range model.Inject {

						model.Inject[x].Name = simpleGenerate("injprop", inj.Name, model)
					}

					o, _ := g.generator.Generate(ctx, &currentTemplate, model, templateOptions)

					code := createCodeResultFile(o, &currentTemplate, model, locationFound, locationPath)
					appendCodeResult(result, code)

				}

				entitiesGenerator()

			} else if currentTemplate.Model == "commits" {
				changelogGenerator := func() {
					cstat.TemplatesUsed += 1

					if ch == nil {
						cp, _ := os.Getwd()
						ch, _ = g.commitHistory.GetCommitHistory(cp, nil)
					}
					model := g.cmc.ToCommitHistoryModel(currentTemplate.Key, currentTemplate.Language, prj, ch)
					// model.PageHeader = simpleGenerate("header", model.PageHeader, model)

					o, _ := g.generator.Generate(ctx, &currentTemplate, model, templateOptions)
					code := createCodeResultFile(o, &currentTemplate, model, locationFound, locationPath)
					appendCodeResult(result, code)

				}

				changelogGenerator()
			} else if currentTemplate.Model == "custom" {
				customGenerator := func() {
					cstat.TemplatesUsed += 1

					if len(options.ModelPath) > 0 {
						dat, err := os.ReadFile(options.ModelPath)

						if err == nil {

							// var customModel interface{}
							err = json.Unmarshal(dat, &options.Custom)
							if err != nil {

							}
						}

					}

					model := g.cmc.ToCustomModel(currentTemplate.Key, currentTemplate.Language, prj, options.Custom)
					// model.PageHeader = simpleGenerate("header", model.PageHeader, model)

					o, _ := g.generator.Generate(ctx, &currentTemplate, model, templateOptions)
					code := createCodeResultFile(o, &currentTemplate, model, locationFound, locationPath)
					appendCodeResult(result, code)

				}

				customGenerator()
			}

		}

		if options.Verbose && !found {
			fmt.Printf("Template '%s' not found \n", tname)
		}

	}
	result.Files = append(result.Files, generatedCode...)

	for _, codeBody := range result.Files {
		if codeBody.Result != nil {

			cstat.Chars += codeBody.Result.Statistics.Chars
			cstat.Lines += codeBody.Result.Statistics.Lines
			cstat.Words += codeBody.Result.Statistics.Words
		}
	}
	cstat.Duration = time.Since(start)

	result.Statistics = cstat
	return result, nil
}

func appendCodeResult(result *models.CodeGenerateResult, code *models.TemplateGeneratorFileResult) {
	if !code.IsSnippet {
		result.Files = append(result.Files, *code)

	} else {
		result.Snippets = append(result.Snippets, *code)
	}
}

func createCodeResultFile(codeResult *models.TemplateGeneratorResult, currentTemplate *models.CodeTemplate, model any, locationFound bool, locationPath string) *models.TemplateGeneratorFileResult {
	fileName, destination := "", ""

	if currentTemplate.FileName != "" {
		fileName = simpleGenerate("filename", currentTemplate.FileName, model)
	}

	if locationFound {
		locationPath = simpleGenerate("location", locationPath, model)
	}

	if len(fileName) > 0 {
		destination = filepath.Join(locationPath, fileName)
	}

	code := models.TemplateGeneratorFileResult{
		Result:      codeResult,
		Destination: destination,
		Filename:    fileName,
		FilePath:    locationPath,
	}

	if currentTemplate.Type == "file" && len(currentTemplate.FileName) > 0 {
		// cstat.FilesCreated += 1
		code.IsSnippet = false
		// result.Files = append(result.Files, code)

	} else if currentTemplate.Type == "snippet" {
		// cstat.SnippetsCreated += 1
		code.IsSnippet = true
		code.SnippetIdentifier = currentTemplate.Identifier
		// result.Snippets = append(result.Snippets, code)
	}

	return &code
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
			for _, templateGroup := range tplVal.Features {

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
func entitiesFromGroups(con *models.ConnectionList, groups []string) []string {
	list := []string{}

	for _, group := range groups {

		conGrp, found := con.Groups[group]
		if found {
			list = append(list, conGrp.Items...)
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
