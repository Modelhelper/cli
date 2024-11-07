package code

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"modelhelper/cli/ports/exporter"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func NewGenerateCodeCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate code",
		Long:  "",
		Run:   codeCommandHandler(app),
	}

	registerFlags(generateCmd)

	return generateCmd
}

func registerFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayP("template", "t", []string{}, "A list of template to convert")
	cmd.Flags().StringArrayP("feature", "f", []string{}, "Use a group of templates")
	// cmd.Flags().StringP("export", "e", "", "Exports the resulted code files to a location or by a mat")

	cmd.Flags().String("template-path", "", "Instructs the program to use this path as root for templates")

	// cmd.Flags().StringP("relations [direct, all, complete]", "r", "", "Include related entities based on the entities in --entity or --entity-group ('direct' | 'all' | 'complete' | 'children' | 'parents')")
	// cmd.Flags().String("template-path", "", "Instructs the program to use this path as root for templates")

	cmd.Flags().StringArrayP("source-group", "g", []string{}, "Use a group of source items (must be defined in the current connection)")
	cmd.Flags().StringArrayP("source", "s", []string{}, "A list of source items to use as a model")

	cmd.Flags().Bool("screen", false, "List the output to the screen, default false")
	cmd.Flags().Bool("copy", false, "Copies the generated code to the clipboard (ctrl + v), default false")
	cmd.Flags().String("export-path", "", "Exports to a directory")
	cmd.Flags().Bool("export-bykey", false, "Exports the code using the template location key, default false")

	cmd.Flags().Bool("overwrite", false, "Overwrite any existing file when exporting to file on disk")
	cmd.Flags().BoolP("verbose", "v", false, "Prints verbose messages, default false")
	cmd.Flags().BoolP("interactive", "i", false, "Goes into interactive mode, default false")

	cmd.Flags().Bool("code-only", false, "Writes only the generated code to the console, no stats, no messages - only code, default false")

	cmd.Flags().Bool("demo", false, "Uses a demo as input source, this will override any other input sources (entity, graphql), default false ")

	// cmd.Flags().String("config-path", "", "Instructs the program to use this config as the config")
	cmd.Flags().String("project-path", "", "Instructs the program to use this project as input")

	// cmd.Flags().String("key", "", "The key to use when encoding and decoding secrets for a connection")

	// cmd.Flags().String("setup", "", "Use this setup to generate code") // version 3.1
	cmd.Flags().StringP("connection", "c", "", "The connection key to be used, uses default connection if not provided")
	cmd.Flags().StringP("name", "n", "", "Sets the name for a template using the 'NameModel'. Will be ignored for any other template type")
	cmd.Flags().String("base-path", "", "Sets the base path for where all templates will be placed in")

	// cmd.Flags().String("custom-json", "", "Instructs the program to use this path as root for templates")
	// cmd.Flags().String("custom-file", "", "Instructs the program to use this path as root for templates")
	cmd.Flags().String("model", "", "Points to a model file to be used as input")

	cmd.RegisterFlagCompletionFunc("relations", completeRelations)
}
func codeCommandHandler(app *modelhelper.ModelhelperCli) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		options := parseCodeOptions(cmd, args)

		if options.RunInteractively {

		}

		result, err := app.Code.Generator.Generate(cmd.Root().Context(), options)
		if err != nil {
			// handle error
			fmt.Println(err)
		}
		sb := strings.Builder{}

		// var fwg sync.WaitGroup
		// var flock = sync.Mutex{}

		for _, res := range result.Files {
			content := []byte(res.Result.Body)
			if options.ExportToScreen {
				screenWriter := exporter.ScreenExporter{}
				screenWriter.Write([]byte(content))
			}

			if options.ExportToClipboard {
				sb.WriteString(string(content))
			}

			if options.ExportByLocationKey && app.Project.Exists {
				fileLocationWriter := exporter.FileExporter{}
				fileLocationWriter.Overwrite = options.Overwrite
				fileLocationWriter.Filename = res.Destination
				_, err := fileLocationWriter.Write(res.Result.Body)
				if err != nil {
					fmt.Printf("Err when writing to '%s', err: %v", res.Destination, err)
				}

			}

			if options.ExportPath != "" {
				// fmt.Printf("Exporting to %s + %s", res.Destination, options.ExportPath)
				fileLocationWriter := exporter.FileExporter{}
				fileLocationWriter.Overwrite = options.Overwrite
				fileLocationWriter.Filename = options.ExportPath
				_, err := fileLocationWriter.Write(res.Result.Body)
				if err != nil {
					fmt.Printf("Err when writing to '%s', err: %v", options.ExportPath, err)
				}
			}
		}

		for _, snippet := range result.Snippets {
			content := []byte(snippet.Result.Body)
			if options.ExportToScreen {
				screenWriter := exporter.ScreenExporter{}
				screenWriter.Write([]byte(content))
			}

			if options.ExportToClipboard {
				sb.WriteString(string(content))
			}

			if options.ExportByLocationKey && app.Project.Exists && len(snippet.Destination) > 0 {
				writer := exporter.NewSnippetExporter(snippet.Destination, snippet.SnippetIdentifier)
				writer.Write(snippet.Result.Body)
			}

		}

		if options.ExportToClipboard {
			fmt.Printf("\nGenerated code is copied to the \033[37mclipboard\033[0m. Use \033[34mctrl+v\033[0m to paste it where you like")
			clipboard.WriteAll(sb.String())
		}

		if !options.CodeOnly {
			printStat(result.Statistics)
		}
		// loc, f := app.Project.ConfigService.FindNearestProjectDir()
		// fmt.Println(loc, f)
	}

}

func parseCodeOptions(cmd *cobra.Command, args []string) *models.CodeGeneratorOptions {
	options := models.CodeGeneratorOptions{}

	codeOnly, _ := cmd.Flags().GetBool("code-only")
	isDemo, _ := cmd.Flags().GetBool("demo")
	sourceFlagArray, _ := cmd.Flags().GetStringArray("source")
	sourceGroupFlagArray, _ := cmd.Flags().GetStringArray("source-group")
	tempPath, _ := cmd.Flags().GetString("template-path")
	projectPath, _ := cmd.Flags().GetString("project-path")
	configFile, _ := cmd.Flags().GetString("config-path")
	inputTemplates, _ := cmd.Flags().GetStringArray("template")
	featureTemplates, _ := cmd.Flags().GetStringArray("feature")
	printScreen, _ := cmd.Flags().GetBool("screen")
	toClipBoard, _ := cmd.Flags().GetBool("copy")
	exportByKey, _ := cmd.Flags().GetBool("export-bykey")
	conName, _ := cmd.Flags().GetString("connection")
	name, _ := cmd.Flags().GetString("name")
	basePathFlag, _ := cmd.Flags().GetString("base-path")
	overwriteAll, _ := cmd.Flags().GetBool("overwrite")
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	interactiveFlag, _ := cmd.Flags().GetBool("interactive")

	options.Name = name
	options.CodeOnly = codeOnly
	options.UseDemo = isDemo
	options.SourceItems = sourceFlagArray
	options.SourceItemGroups = sourceGroupFlagArray
	options.TemplatePath = tempPath
	options.ConfigFilePath = configFile
	options.ProjectFilePath = projectPath
	options.Templates = inputTemplates
	options.FeatureTemplates = featureTemplates
	options.ExportToScreen = printScreen
	options.ExportToClipboard = toClipBoard
	options.ExportByLocationKey = exportByKey
	options.ConnectionName = conName
	options.Overwrite = overwriteAll
	options.Verbose = verboseFlag
	options.RunInteractively = interactiveFlag
	options.ModelPath, _ = cmd.Flags().GetString("model")
	options.ExportPath, _ = cmd.Flags().GetString("export-path")
	if len(basePathFlag) > 0 {
		options.BasePath = basePathFlag
	}

	options.CanUseTemplates = len(options.Templates) > 0 || len(options.FeatureTemplates) > 0
	return &options
}

func completeRelations(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"direct", "all", "complete", "children", "parents"}, cobra.ShellCompDirectiveDefault
}

func exportPath() {

}

func printStat(stat *models.TemplateGeneratorStatistics) {
	if stat == nil {
		return
	}
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

	wpm := 30.0
	cpm := 250.0

	min := float64(stat.Words) / wpm
	min = float64(stat.Chars) / cpm

	fmt.Printf("\nIn summary... It took \033[32m%vms\033[0m to generate \033[34m%d\033[0m words and \033[34m%d\033[0m lines. \nYou saved around \033[32m%v minutes\033[0m by not typing it youreself\n",
		stat.Duration.Milliseconds(),
		stat.Words,
		stat.Lines,
		int(min))
}
