package cmd

import (
	"context"

	"fmt"
	"log"
	"modelhelper/cli/app"
	"modelhelper/cli/code"
	"modelhelper/cli/codegen"
	"modelhelper/cli/config"
	"modelhelper/cli/modelhelper"
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

type codeGenerateOptions struct {
	templates         []string
	templateGroups    []string
	templatePath      string
	canUseTemplates   bool
	entityGroups      []string
	entities          []string
	exportToScreen    bool
	exportByKey       bool
	exportPath        string
	connection        string
	exportToClipboard bool
	overwrite         bool
	relations         string
	codeOnly          bool
	useDemo           bool
	configFilePath    string
	projectFilePath   string
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
	// generateCmd.Flags().Bool("export-bykey", false, "Exports the code using the template keys, default false")
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

func newGenerateOptions(cmd *cobra.Command, args []string) *codeGenerateOptions {
	options := codeGenerateOptions{}

	codeOnly, _ := cmd.Flags().GetBool("code-only")
	isDemo, _ := cmd.Flags().GetBool("demo")
	entityFlagArray, _ := cmd.Flags().GetStringArray("entity")
	entityGroupFlagArray, _ := cmd.Flags().GetStringArray("entity-group")
	tempPath, _ := cmd.Flags().GetString("template-path")
	projectPath, _ := cmd.Flags().GetString("project-path")
	configFile, _ := cmd.Flags().GetString("config-path")
	exportFile, _ := cmd.Flags().GetString("export-path")
	inputTemplates, _ := cmd.Flags().GetStringArray("template")
	inputGroupTemplates, _ := cmd.Flags().GetStringArray("template-group")
	printScreen, _ := cmd.Flags().GetBool("screen")
	toClipBoard, _ := cmd.Flags().GetBool("copy")
	// exportByKey, _ := cmd.Flags().GetBool("export-bykey")
	conName, _ := cmd.Flags().GetString("connection")
	overwriteAll, _ := cmd.Flags().GetBool("overwrite")

	options.codeOnly = codeOnly
	options.useDemo = isDemo
	options.entities = entityFlagArray
	options.entityGroups = entityGroupFlagArray
	options.templatePath = tempPath
	options.configFilePath = configFile
	options.projectFilePath = projectPath
	options.templates = inputTemplates
	options.templateGroups = inputGroupTemplates
	options.exportToScreen = printScreen
	options.exportToClipboard = toClipBoard
	options.exportByKey = false
	options.exportPath = exportFile
	options.connection = conName
	options.overwrite = overwriteAll

	options.canUseTemplates = len(options.templates) > 0 || len(options.templateGroups) > 0
	return &options
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g", "gen"},
	Short:   "Generates code based on language, template and source",

	Run: func(cmd *cobra.Command, args []string) {

		options := newGenerateOptions(cmd, args)

		if len(options.templates) == 0 && len(options.templateGroups) == 0 {
			// no point to continue if no templates is given
			fmt.Printf(`No templates or template groups are provided resulting in nothing to create
please use mh generate with the -t or --template [templatename] to set at template

You could also use mh template or mh t to see a list of all available templates`)

			return
		}

		// obsolete
		modelHelperApp = app.New()

		// obsolete
		appCtx := modelHelperApp.CreateContext()

		var con modelhelper.Connection
		var prj *modelhelper.ProjectConfig
		var entities []modelhelper.Entity

		cfg, err := loadConfig(options.configFilePath)
		if err != nil {
			// handle error
		}

		if options.useDemo {
			options.connection = "demo"
			con = modelhelper.Connection{Type: options.connection}
		} else {
			if len(appCtx.Connections) == 0 {
				fmt.Println("Could not find any connections to use, please add a connection")
				fmt.Println("to the config and/or any project file")
				return
			}
			if len(options.connection) == 0 {

				options.connection = appCtx.DefaultConnection
			}

			if len(options.connection) == 0 {
				ka := keyArray(appCtx.Connections)
				options.connection = ka[0]
			}

			con = appCtx.Connections[options.connection]
		}

		entityList := mergedList(options.entities, entitiesFromGroups(con, options.entityGroups))

		src := source.SourceFactory(&con)

		entities = *loadEntities(src, entityList)

		prj = loadProject(options.projectFilePath)

		languages, _ := code.LoadFromPath(cfg.Languages.Definitions)

		if len(options.templatePath) == 0 {
			options.templatePath = cfg.Templates.Location
		}

		allTemplates, blocks := loadTemplates(options.templatePath)

		// if err != nil {
		// 	panic(err)
		// }

		options.templates = selectTemplates(allTemplates, options.templates, options.templateGroups)

		start := time.Now()
		var cstat = modelhelper.CodeGeneratorStatistics{}
		var generatedCode []modelhelper.CodeFileResult

		// creates the root context to be passed to each sub routine
		ctxVal := codegen.CodeContextValue{}
		ctxVal.Blocks = blocks

		for _, tname := range options.templates {

			// var tt *tpl.Template
			// fmt.Println(tname)
			currentTemplate, found := allTemplates[tname]

			if found {
				var codeSection modelhelper.Code
				// var csFound = false
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
					// codeSection, csFound = prj.Code[currentTemplate.Language]
					codeSection = prj.Code[currentTemplate.Language]
				}
				// }
				generator := codegen.GoLangGenerator{}

				ctx := context.WithValue(context.Background(), "code", ctxVal)
				if len(currentTemplate.Model) == 0 || currentTemplate.Model == "basic" {

					basicGenerator := func() {
						cstat.TemplatesUsed += 1

						model := ToBasicModel(currentTemplate.Key, currentTemplate.Language, prj)
						o, _ := generator.Generate(ctx, model)

						f := modelhelper.CodeFileResult{
							Result:   o,
							Filename: "",
						}

						generatedCode = append(generatedCode, f)
					}

					basicGenerator()

				} else if currentTemplate.Model == "entity" && len(entities) > 0 {

					for _, entity := range entities {

						entityGenerator := func() {
							cstat.TemplatesUsed += 1
							cstat.EntitiesUsed += 1

							model := ToEntityModel(currentTemplate.Key, currentTemplate.Language, prj, &entity)

							model.PageHeader = codegen.Generate("header", model.PageHeader, model)
							model.Namespace = codegen.Generate("namesp", model.Namespace, model)

							for i, imp := range model.Imports {

								model.Imports[i] = codegen.Generate("import", imp, model)
							}

							model.Imports = removeDuplicateStringValues(model.Imports)

							for x, inj := range model.Inject {

								model.Inject[x].Name = codegen.Generate("injprop", inj.Name, model)
							}

							o, _ := generator.Generate(ctx, model)

							// fullPath := ""
							fileName := ""
							if currentTemplate.Type == "file" && len(currentTemplate.FileName) > 0 {
								cstat.FilesCreated += 1

								fileName = codegen.Generate("filename", currentTemplate.FileName, model)

								// if csFound {
								// 	if options.exportPath
								// 	fullPath = filepath.Join(codeSection.Locations[currentTemplate.Key], filen)
								// }
							}

							f := modelhelper.CodeFileResult{
								Result:   o,
								Filename: fileName,
								FilePath: codeSection.Locations[currentTemplate.Key],
							}

							generatedCode = append(generatedCode, f)
						}

						entityGenerator()

					}
				} else if currentTemplate.Model == "entities" && len(entities) > 0 {

					entitiesGenerator := func() {
						cstat.TemplatesUsed += 1
						model := ToEntitiesModel(currentTemplate.Key, currentTemplate.Language, prj, &entities)
						model.PageHeader = codegen.Generate("header", model.PageHeader, model)

						model.Namespace = codegen.Generate("namesp", model.Namespace, model)

						for i, imp := range model.Imports {

							model.Imports[i] = codegen.Generate("import", imp, model)
						}

						model.Imports = removeDuplicateStringValues(model.Imports)

						for x, inj := range model.Inject {

							model.Inject[x].Name = codegen.Generate("injprop", inj.Name, model)
						}

						o, _ := generator.Generate(ctx, model)

						fileName := ""
						if currentTemplate.Type == "file" && len(currentTemplate.FileName) > 0 {
							cstat.FilesCreated += 1
							fileName = codegen.Generate("filename", currentTemplate.FileName, model)
							// if csFound {

							// 	fullPath = filepath.Join(codeSection.Locations[currentTemplate.Key], filen)
							// }
						}

						f := modelhelper.CodeFileResult{
							Result:   o,
							Filename: fileName,
							FilePath: codeSection.Locations[currentTemplate.Key],
						}

						generatedCode = append(generatedCode, f)

					}

					entitiesGenerator()

				}

			}

		}

		sb := strings.Builder{}
		var fwg sync.WaitGroup
		var flock = sync.Mutex{}
		for _, codeBody := range generatedCode {
			cstat.Chars += codeBody.Result.Statistics.Chars
			cstat.Lines += codeBody.Result.Statistics.Lines
			cstat.Words += codeBody.Result.Statistics.Words
			content := []byte(codeBody.Result.Body)
			if options.exportToScreen {
				screenWriter := tpl.ScreenExporter{}
				screenWriter.Write([]byte(content))
			}

			if options.exportToClipboard {
				sb.WriteString(string(codeBody.Result.Body))
			}

			if options.exportByKey && len(codeBody.Filename) > 0 {
				fwg.Add(1)
				go func(filename string, rootPath string, content []byte) {
					defer fwg.Done()
					keyExporter := tpl.FileExporter{
						Filename:  filepath.Join(rootPath, filename),
						Overwrite: options.overwrite,
					}

					_, err := keyExporter.Write([]byte(content))
					if err != nil {
						fmt.Println(filepath.ErrBadPattern)
					}
					// fmt.Println("*** FILENAME::", s.filename)
					flock.Lock()
					cstat.FilesExported += 1
					flock.Unlock()
				}(codeBody.Filename, "D:/projects/ModelHelper", content)

			}
			// TODO: export to file
			if len(options.exportPath) > 0 {
				fwg.Add(1)
				go func(filename string, rootPath string, content []byte) {
					defer fwg.Done()

					if len(filename) > 0 {

						fileExporter := tpl.FileExporter{
							Filename:  filepath.Join(rootPath, filename),
							Overwrite: options.overwrite,
						}

						_, err := fileExporter.Write([]byte(content))
						if err != nil {
							fmt.Printf("%s, err: \n%v", filepath.ErrBadPattern, err)
						}
						// fmt.Println("*** FILENAME::", s.filename)
						flock.Lock()
						cstat.FilesExported += 1
						flock.Unlock()
					} else {
						fmt.Println("Filename empty...")
					}
				}(codeBody.Filename, options.exportPath, content)
			}
		}

		fwg.Wait()
		if options.exportToClipboard {
			fmt.Printf("\nGenerated code is copied to the \033[37mclipboard\033[0m. Use \033[34mctrl+v\033[0m to paste it where you like")
			clipboard.WriteAll(sb.String())
		}

		cstat.Duration = time.Since(start)
		// stat["total.time"] = int(cstat.duration.Milliseconds())
		if !options.codeOnly {
			wpm := 30.0
			cpm := 250.0

			min := float64(cstat.Words) / wpm
			min = float64(cstat.Chars) / cpm
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

func selectTemplates(templates map[string]tpl.Template, input []string, groups []string) []string {
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
func entitiesFromGroups(con modelhelper.Connection, groups []string) []string {
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

// type codeFile struct {
// 	filename        string
// 	filePath        string
// 	result          codegen.Result
// 	exists          bool
// 	existingContent string
// }

func completeRelations(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"direct", "all", "complete", "children", "parents"}, cobra.ShellCompDirectiveDefault
}
func printStat(stat modelhelper.CodeGeneratorStatistics) {
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

type basicModel struct {
	project  *modelhelper.ProjectConfig
	key      string
	language string
}

type entityModel struct {
	entity   *modelhelper.Entity
	project  *modelhelper.ProjectConfig
	key      string
	language string
}

type entitiesModel struct {
	entities *[]modelhelper.Entity
	project  *modelhelper.ProjectConfig
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

func loadEntities(src source.Source, names []string) *[]modelhelper.Entity {
	var entities []modelhelper.Entity

	// if isDemo {
	// 	// load demo project
	// 	var ds *source.DemoSource
	// 	el, _ := ds.Entities("")

	// 	for _, eitem := range *el {
	// 		entities = append(entities, eitem)
	// 	}

	// 	// load demo tables (2)
	// } else {

	// conName := appCtx.DefaultConnection

	// con := appCtx.Connections[conName]
	// src := con.LoadSource()

	if len(names) > 0 {
		var wg sync.WaitGroup
		var lock = sync.Mutex{}
		for _, entityName := range names {

			wg.Add(1)
			go func(name string) {
				defer wg.Done()
				entity, err := src.Entity(name)
				if err != nil {
					log.Fatalln(err)
				}

				lock.Lock()
				entities = append(entities, *entity)
				lock.Unlock()
			}(entityName)
		}

		wg.Wait()
	}

	// }

	return &entities
}

func loadConfig(configPath string) (*modelhelper.Config, error) {
	cfg := config.NewConfigLoader()

	if len(configPath) > 0 {
		return cfg.LoadFromFile(configPath)
	} else {
		return cfg.Load()
	}
}
func loadProject(projectPath string) *modelhelper.ProjectConfig {
	p := project.NewModelhelperProject()

	if len(projectPath) == 0 {
		if project.Exists(project.DefaultLocation()) {
			projectPath = project.DefaultLocation()
		} else {
			fp, foundProject := p.FindNearestProjectDir()
			if foundProject && project.Exists(fp) {
				projectPath = fp

			}
		}
	}

	if len(projectPath) > 0 {

		prj, _ := project.Load(projectPath)
		return &prj.Config
	}

	return nil
}

func toProjectModel(input modelhelper.ProjectConfig) modelhelper.ProjectSection {
	return modelhelper.ProjectSection{
		Owner: input.OwnerName,
		Name:  input.Name,
	}

}

func ToEntityModel(key, language string, project *modelhelper.ProjectConfig, entity *modelhelper.Entity) modelhelper.EntityModel {

	entityBase := modelhelper.EntityModel{}

	base := ToBasicModel(key, language, project)
	if entity != nil {
		entityBase = toEntitySection(entity)
	}

	out := modelhelper.EntityModel{
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
		HasSynonym:                entityBase.HasSynonym,
		Synonym:                   entityBase.Synonym,
		ModelName:                 entityBase.ModelName,
	}

	out.HasChildren = len(entityBase.Children) > 0
	out.HasParents = len(entityBase.Parents) > 0
	return out
}
func ToEntitiesModel(key, language string, project *modelhelper.ProjectConfig, entities *[]modelhelper.Entity) modelhelper.EntityListModel {

	entitylist := []modelhelper.EntityModel{}

	base := ToBasicModel(key, language, project)
	if entities != nil && len(*entities) > 0 {
		for _, entity := range *entities {

			entityBase := toEntitySection(&entity)
			entitylist = append(entitylist, entityBase)
		}
	}

	out := modelhelper.EntityListModel{
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

func ToBasicModel(key, language string, project *modelhelper.ProjectConfig) modelhelper.BasicModel {
	b := modelhelper.BasicModel{}
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

				b.Inject = []modelhelper.InjectSection{}
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

func toColumnSection(from modelhelper.Column, entityName string) modelhelper.EntityColumnModel {
	col := modelhelper.EntityColumnModel{
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
func toEntitySection(from *modelhelper.Entity) modelhelper.EntityModel {
	out := modelhelper.EntityModel{
		Name:               from.Name,
		Schema:             from.Schema,
		Type:               from.Type,
		Alias:              from.Alias,
		Description:        from.Description,
		HasDescription:     len(from.Description) > 0,
		HasPrefix:          false,
		NameWithoutPrefix:  "",
		UsesIdentityColumn: from.UsesIdentityColumn,
		HasSynonym:         from.HasSynonym,
	}

	out.Synonym = from.Synonym

	out.ModelName = coalesceString(from.Synonym, from.Name)

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
		child := modelhelper.EntityRelationModel{}

		child.Name = cr.Name
		child.Schema = cr.Schema
		child.RelatedColumn = modelhelper.EntityColumnProps{
			Name:       cr.ColumnName,
			DataType:   cr.ColumnType,
			IsNullable: cr.ColumnNullable,
		}

		child.OwnerColumn = modelhelper.EntityColumnProps{
			Name:       cr.OwnerColumnName,
			DataType:   cr.OwnerColumnType,
			IsNullable: cr.OwnerColumnNullable,
		}

		out.Children = append(out.Children, child)
		// children = append(children)
		child.NameWithoutPrefix = strings.TrimPrefix(cr.Name, out.Name)
		child.HasPrefix = strings.HasPrefix(cr.Name, out.Name)

		child.HasSynonym = cr.HasSynonym
		if child.HasSynonym {
			child.Synonym = cr.Synonym
		}

		child.ModelName = coalesceString(cr.Synonym, cr.Name)
	}

	for _, pr := range from.ParentRelations {
		parent := modelhelper.EntityRelationModel{}
		parent.Name = pr.Name
		parent.Schema = pr.Schema

		parent.HasDescription = false

		parent.HasSynonym = pr.HasSynonym
		if parent.HasSynonym {
			parent.Synonym = pr.Synonym
		}

		parent.OwnerColumn = modelhelper.EntityColumnProps{
			Name:       pr.ColumnName,
			DataType:   pr.ColumnType,
			IsNullable: pr.ColumnNullable,
		}

		parent.RelatedColumn = modelhelper.EntityColumnProps{
			Name:       pr.OwnerColumnName,
			DataType:   pr.OwnerColumnType,
			IsNullable: pr.OwnerColumnNullable,
		}

		parent.ModelName = coalesceString(pr.Synonym, pr.Name)

		parent.NameWithoutPrefix = strings.TrimPrefix(pr.Name, out.Name)
		parent.HasPrefix = strings.HasPrefix(pr.Name, out.Name)

		out.Parents = append(out.Parents, parent)
	}

	return out
}
func toInjectSection(from modelhelper.Inject, m interface{}) modelhelper.InjectSection {
	// name, _ := codegen.Generate("fileName", from.Name, m)
	code := modelhelper.InjectSection{
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
			Creator:  tpl.CreatorImportModel{CompanyName: "FooBar inc", DeveloperName: "Dev E. Loper"},
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

func coalesceString(name ...string) string {
	output := ""
	for _, n := range name {
		if len(n) > 0 {
			output = n
			break
		}
	}

	return output
}
