package code

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"

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
	cmd.Flags().StringArray("template-group", []string{}, "Use a group of templates")
	cmd.Flags().String("template-path", "", "Instructs the program to use this path as root for templates")

	cmd.Flags().StringP("relations [direct, all, complete]", "r", "", "Include related entities based on the entities in --entity or --entity-group ('direct' | 'all' | 'complete' | 'children' | 'parents')")
	// cmd.Flags().String("template-path", "", "Instructs the program to use this path as root for templates")

	cmd.Flags().StringArray("entity-group", []string{}, "Use a group of entities (must be defines in the current connection)")
	cmd.Flags().StringArrayP("entity", "e", []string{}, "A list of entits to use as a model")

	cmd.Flags().Bool("screen", false, "List the output to the screen, default false")
	cmd.Flags().Bool("copy", false, "Copies the generated code to the clipboard (ctrl + v), default false")
	cmd.Flags().String("export-path", "", "Exports to a directory")
	// cmd.Flags().Bool("export-bykey", false, "Exports the code using the template keys, default false")
	cmd.Flags().Bool("overwrite", false, "Overwrite any existing file when exporting to file on disk")

	cmd.Flags().Bool("code-only", false, "Writes only the generated code to the console, no stats, no messages - only code, default false")

	cmd.Flags().Bool("demo", false, "Uses a demo as input source, this will override any other input sources (entity, graphql), default false ")

	cmd.Flags().String("config-path", "", "Instructs the program to use this config as the config")
	cmd.Flags().String("project-path", "", "Instructs the program to use this project as input")

	cmd.Flags().String("key", "", "The key to use when encoding and decoding secrets for a connection")

	// cmd.Flags().String("setup", "", "Use this setup to generate code") // version 3.1
	cmd.Flags().StringP("connection", "c", "", "The connection key to be used, uses default connection if not provided")

	cmd.RegisterFlagCompletionFunc("relations", completeRelations)
}
func codeCommandHandler(app *modelhelper.ModelhelperCli) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		options := parseCodeOptions(cmd, args)
		result, err := app.Code.Generator.Generate(cmd.Root().Context(), options)

		if err != nil {
			// handle error
			fmt.Println(err)
		}

		for _, res := range result {
			fmt.Printf("Printing the generated result:\n%s", string(res.Result.Body))
		}
	}
}

func parseCodeOptions(cmd *cobra.Command, args []string) *models.CodeGeneratorOptions {
	options := models.CodeGeneratorOptions{}

	codeOnly, _ := cmd.Flags().GetBool("code-only")
	isDemo, _ := cmd.Flags().GetBool("demo")
	entityFlagArray, _ := cmd.Flags().GetStringArray("entity")
	entityGroupFlagArray, _ := cmd.Flags().GetStringArray("entity-group")
	tempPath, _ := cmd.Flags().GetString("template-path")
	projectPath, _ := cmd.Flags().GetString("project-path")
	configFile, _ := cmd.Flags().GetString("config-path")
	inputTemplates, _ := cmd.Flags().GetStringArray("template")
	inputGroupTemplates, _ := cmd.Flags().GetStringArray("template-group")
	printScreen, _ := cmd.Flags().GetBool("screen")
	toClipBoard, _ := cmd.Flags().GetBool("copy")
	// exportByKey, _ := cmd.Flags().GetBool("export-bykey")
	conName, _ := cmd.Flags().GetString("connection")
	overwriteAll, _ := cmd.Flags().GetBool("overwrite")

	options.CodeOnly = codeOnly
	options.UseDemo = isDemo
	options.Entities = entityFlagArray
	options.EntityGroups = entityGroupFlagArray
	options.TemplatePath = tempPath
	options.ConfigFilePath = configFile
	options.ProjectFilePath = projectPath
	options.Templates = inputTemplates
	options.TemplateGroups = inputGroupTemplates
	options.ExportToScreen = printScreen
	options.ExportToClipboard = toClipBoard
	options.ExportByKey = false
	options.ConnectionName = conName
	options.Overwrite = overwriteAll

	options.CanUseTemplates = len(options.Templates) > 0 || len(options.TemplateGroups) > 0
	return &options
}

func completeRelations(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"direct", "all", "complete", "children", "parents"}, cobra.ShellCompDirectiveDefault
}
