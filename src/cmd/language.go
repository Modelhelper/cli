package cmd

import (
	"fmt"
	"modelhelper/cli/code"
	"modelhelper/cli/config"
	"modelhelper/cli/ui"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// aboutCmd represents the about command
var languageCmd = &cobra.Command{
	Use:     "language",
	Aliases: []string{"lang", "l"},
	Short:   "Show a list of language definitions installed for modelhelper",
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		openFile, _ := cmd.Flags().GetBool("open")

		defs, err := code.LoadFromPath(cfg.Languages.Definitions)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		if len(args) > 0 {
			def := defs[args[0]]

			if openFile {
				openPathInEditor(cfg.DefaultEditor, def.Path)
			} else {
				presentLanguage(&def)

			}

			// showLanguage(args[0], cfg.Languages.Definitions)
		} else {
			ui.ConsoleTitle("Language list")
			fmt.Println(`
This is a list of all available languages defined for model helper			
			`)
			renderer := languageTableRenderer{defs}
			ui.RenderTable(&renderer, &renderer)
		}

	},
}

func presentLanguage(def *code.LanguageDefinition) {

	ui.PrintConsoleTitle(fmt.Sprintf("Language Definition for %s", def.Language))

	if len(def.Description) > 0 {
		fmt.Println(def.Description)
	} else {
		fmt.Println(def.Short)
	}
	fmt.Printf("\n%-20s%8s\n", "Language:", def.Language)
	fmt.Printf("%-20s%8s\n\n", "Version:", def.Version)
	fmt.Printf("%-20s%8d\n", "Datatypes defined:", len(def.DataTypes))
	fmt.Printf("%-20s%8d\n", "Keys defined:", len(def.Keys))
	fmt.Printf("%-20s%8d\n\n", "Injects defined:", len(def.Inject))

	if def.DataTypes != nil {
		ui.PrintConsoleTitle("Datatypes")
		dtr := datatypeTableRenderer{def.Language, def.DataTypes}
		ui.RenderTable(&dtr, &dtr)
	}

	if def.Keys != nil {

		ui.PrintConsoleTitle("Keys")
		kr := keyRenderer{def.Keys}
		ui.RenderTable(&kr, &kr)
	}

	if def.Inject != nil {
		ui.PrintConsoleTitle("Injects")

		ir := injectRenderer{def.Inject}
		ui.RenderTable(&ir, &ir)
	}

	fmt.Println("")
}

func init() {
	rootCmd.AddCommand(languageCmd)
	languageCmd.Flags().Bool("open", false, "Opens the language definition in the default editor")
}

type languageTableRenderer struct {
	rows map[string]code.LanguageDefinition
}

func (d *languageTableRenderer) BuildHeader() []string {
	h := []string{"Language", "Version", "Datatypes", "Imports", "Keys", "Injects", "Description"}

	return h
}
func (d *languageTableRenderer) ToRows() [][]string {
	var rows [][]string

	p := message.NewPrinter(language.English)

	for _, row := range d.rows {
		// un := "No"
		// ci := "No"
		r := []string{
			row.Language,
			row.Version,
			p.Sprintf("%d", len(row.DataTypes)),
			p.Sprintf("%d", len(row.DefaultImports)),
			p.Sprintf("%d", len(row.Keys)),
			p.Sprintf("%d", len(row.Inject)),
			row.Short,
		}

		rows = append(rows, r)
	}

	return rows
}

type datatypeTableRenderer struct {
	language  string
	datatypes map[string]code.Datatype
}

func (l *datatypeTableRenderer) BuildHeader() []string {
	return []string{
		"DB Datatype",
		fmt.Sprintf("%s not null", l.language),
		fmt.Sprintf("%s nullable", l.language),
		fmt.Sprintf("%s nullable(alt)", l.language),
	}
}

func (d *datatypeTableRenderer) ToRows() [][]string {
	var rows [][]string

	for key, val := range d.datatypes {
		r := []string{
			key,
			val.NotNull,
			val.Nullable,
			val.NullableAlternative,
		}

		rows = append(rows, r)
	}

	return rows
}

type keyRenderer struct {
	keys map[string]code.Key
}

func (l *keyRenderer) BuildHeader() []string {
	return []string{
		"Key",
		"Namespace",
		"Prefix",
		"Postfix",
		"Import count",
		"Injects",
	}
}

func (d *keyRenderer) ToRows() [][]string {
	var rows [][]string
	p := message.NewPrinter(language.English)

	for key, val := range d.keys {
		// injArray := []string{}
		// for k, _ := range val.Inject {
		// 	injArray = append(injArray, k)
		// }
		inj := strings.Join(val.Inject, ", ")
		r := []string{
			key,
			val.Namespace,
			val.Prefix,
			val.Postfix,
			p.Sprintf("%d", len(val.Imports)),
			inj,
		}

		rows = append(rows, r)
	}

	return rows
}

type injectRenderer struct {
	keys map[string]code.Inject
}

func (l *injectRenderer) BuildHeader() []string {
	return []string{
		"Key",
		"Name",
		"Property",
		"Import count",
	}
}

func (d *injectRenderer) ToRows() [][]string {
	var rows [][]string
	p := message.NewPrinter(language.English)

	for key, val := range d.keys {

		r := []string{
			key,
			val.Name,
			val.PropertyName,
			p.Sprintf("%d", len(val.Imports)),
		}

		rows = append(rows, r)
	}

	return rows
}
