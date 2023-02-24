package project

import (
	"fmt"
	"io/ioutil"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"os"
	"path/filepath"
	"time"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

type generateOptions struct {
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

func NewGenerateProjectCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	generateCmd := &cobra.Command{
		Use: "create",
		// Args: cobra.RangeArgs(1, 2),
		Args:  cobra.ExactArgs(0),
		Short: "Creates a new project based on a template",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

			start := time.Now()
			fileCnt := 0
			isDryRun, _ := cmd.Flags().GetBool("dry-run")
			opt, dest := parseProjectArgs(cmd, args)

			if len(opt.Name) > 0 && len(opt.Template) > 0 {
				tpl := app.Project.TemplateService.Load(opt.Template)

				if tpl != nil {
					dr := ""
					if isDryRun {
						dr = color.Red.Sprintf("\n** This is a Dry Run, no files will be written to the system **\n\n")
					}
					fmt.Printf("%sTemplate: \n\tname: %s\n\tlang: %s\n", dr, tpl.Name, tpl.Language)

					mdl := app.Project.Generator.BuildTemplateModel(opt, tpl)
					src, err := app.Project.Generator.Generate(cmd.Root().Context(), tpl, mdl)
					if err != nil {
						fmt.Printf("Err %v", err)
					}

					fileCnt = len(src)

					root, err := app.Project.Generator.GenerateRootDirectoryName(tpl.RootDirectory, mdl)
					if err != nil {
						fmt.Printf("Err %v", err)
					}

					path := writeLocation(dest, root)

					if !isDryRun {
						writeFilesToLocation(path, src)
					}

					fmt.Printf("\nSource files created at this location: \n\t** %s **\n\n", path)
					for _, srcFile := range src {
						fmt.Printf("\t%s\n", filepath.Join(path, srcFile.FileName))
					}

				}
			} else {
				fmt.Println("No input for name or template")
			}

			elapsed := time.Since(start)

			fmt.Printf("\n\nGenerating %d files for project '%s' took %dms\n", fileCnt, opt.Name, elapsed.Milliseconds())

		},
	}

	generateCmd.Flags().StringP("name", "n", "", "Name of the project")
	generateCmd.Flags().StringP("template", "t", "", "The template to use")
	generateCmd.Flags().StringP("destination", "d", "", "Exports to a directory")
	generateCmd.Flags().Bool("dry-run", false, "If set to true, no files will be written to the system")

	// generateCmd.Flags().StringArray("template-group", []string{}, "Use a group of templates")
	// generateCmd.Flags().String("template-path", "", "Instructs the program to use this path as root for templates")

	// generateCmd.Flags().StringP("relations [direct, all, complete]", "r", "", "Include related entities based on the entities in --entity or --entity-group ('direct' | 'all' | 'complete' | 'children' | 'parents')")
	// generateCmd.Flags().String("template-path", "", "Instructs the program to use this path as root for templates")

	// generateCmd.Flags().StringArray("entity-group", []string{}, "Use a group of entities (must be defines in the current connection)")
	// generateCmd.Flags().StringArrayP("entity", "e", []string{}, "A list of entits to use as a model")

	// generateCmd.Flags().Bool("screen", false, "List the output to the screen, default false")
	// generateCmd.Flags().Bool("copy", false, "Copies the generated code to the clipboard (ctrl + v), default false")
	// generateCmd.Flags().Bool("export-bykey", false, "Exports the code using the template keys, default false")
	// generateCmd.Flags().Bool("overwrite", false, "Overwrite any existing file when exporting to file on disk")

	// generateCmd.Flags().Bool("code-only", false, "Writes only the generated code to the console, no stats, no messages - only code, default false")

	// generateCmd.Flags().String("config-path", "", "Instructs the program to use this config as the config")
	// generateCmd.Flags().String("project-path", "", "Instructs the program to use this project as input")

	// generateCmd.Flags().String("key", "", "The key to use when encoding and decoding secrets for a connection")

	// generateCmd.Flags().String("setup", "", "Use this setup to generate code") // version 3.1
	// generateCmd.Flags().StringP("connection", "c", "", "The connection key to be used, uses default connection if not provided")

	// generateCmd.RegisterFlagCompletionFunc("relations", completeRelations)

	return generateCmd
}

func parseProjectArgs(cmd *cobra.Command, args []string) (*models.ProjectTemplateCreateOptions, string) {
	options := models.ProjectTemplateCreateOptions{}

	template, _ := cmd.Flags().GetString("template")
	name, _ := cmd.Flags().GetString("name")
	dest, _ := cmd.Flags().GetString("destination")

	options.Name = name
	options.Template = template
	// options.Destination = d

	return &options, dest
}

func completeRelations(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"direct", "all", "complete", "children", "parents"}, cobra.ShellCompDirectiveDefault
}

func writeLocation(dest, rootFolder string) string {
	if dest == "" {
		dest = currentDirectory()
	}

	return filepath.Join(dest, rootFolder)
}
func currentDirectory() string {
	cd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error when getting currentDirectory %v", err)
		panic(err)
	}

	return cd
}

func writeFilesToLocation(destination string, files []*models.ProjectSourceFile) {
	for _, file := range files {
		path := filepath.Join(destination, file.RelativePath, file.FileName)

		d, _ := filepath.Split(path)
		if _, err := os.Stat(d); os.IsNotExist(err) {
			err := os.MkdirAll(d, 0700)
			if err != nil {
				fmt.Printf("Could not create '%s'", d)
			}
		}

		err := ioutil.WriteFile(path, file.Content, 0777)
		if err != nil {
			fmt.Printf("Could not write '%s' to disk", path)
		}
	}
}
