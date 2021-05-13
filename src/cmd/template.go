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
	"modelhelper/cli/app"
	"modelhelper/cli/config"
	"modelhelper/cli/tpl"
	"modelhelper/cli/ui"
	"strings"

	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "List all templates or view content of a single template",
	Long: `With this command you can list all available templates, snippet and blocks

with the use of --by option you can group the various templates either by type, language
group or tag.

A template can exist in one or more groups making it possible to generate different
outcome based on what you need

Filtering the templates:
Filter the template by using on or more of the following options
--lang <langcodes> (e.g --lang cs), filters by language
--type <type> (e.g --type block), filters by type
--model <model> (e.g --model entity), filters by model
--group <groupname> (e.g --group cs-dpr-full), filters by group


-- hva med liste på gruppenavn
-- hva med liste på tags
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()

		tl := tpl.TemplateLoader{
			Directory: app.TemplateFolder(cfg.Templates.Location),
		}
		allTemplates, _ := tl.LoadTemplates()

		group, _ := cmd.Flags().GetString("by")
		typeFiler, _ := cmd.Flags().GetStringArray("type")
		langFilter, _ := cmd.Flags().GetStringArray("lang")
		modelFilter, _ := cmd.Flags().GetStringArray("model")
		keyFilter, _ := cmd.Flags().GetStringArray("key")

		var tp templatePrinter
		if len(typeFiler) > 0 {
			ft := filterByType{}
			allTemplates = ft.filter(allTemplates, typeFiler)
		}

		if len(langFilter) > 0 {
			ft := filterByLang{}
			allTemplates = ft.filter(allTemplates, langFilter)
		}
		if len(keyFilter) > 0 {
			ft := filterByKey{}
			allTemplates = ft.filter(allTemplates, keyFilter)
		}
		if len(modelFilter) > 0 {
			ft := filterByModel{}
			allTemplates = ft.filter(allTemplates, modelFilter)
		}

		if len(group) > 0 {
			ui.PrintConsoleTitle("ModelHelper Templates grouped by " + group)
			fmt.Println(`In the list below you will find all available templates in ModelHelper`)
			grouper := getGrouper(strings.ToLower(group))
			mg := grouper.group(allTemplates)

			for typ, tv := range mg {
				ui.ConsoleTitle(typ)

				tp.templates = tv
				ui.RenderTable(&tp, &tp)

			}
		} else {
			ui.PrintConsoleTitle("ModelHelper Templates")
			fmt.Println(`In the list below you will find all available templates in ModelHelper`)

			tp.templates = allTemplates
			ui.RenderTable(&tp, &tp)
		}

	},
}

func init() {

	rootCmd.AddCommand(templateCmd)
	templateCmd.Flags().String("by", "", "Groups the templates by type, group, language, model or tag")
	templateCmd.Flags().StringArray("type", []string{}, "Filter the templates by the name of the type")
	templateCmd.Flags().StringArray("lang", []string{}, "Filter the templates by language")
	templateCmd.Flags().StringArray("model", []string{}, "Filter the templates by model")
	templateCmd.Flags().StringArray("key", []string{}, "Filter the templates by model")

}

type templatePrinter struct {
	templates map[string]tpl.Template
}

func (t *templatePrinter) ToRows() [][]string {
	var rows [][]string

	for name, t := range t.templates {
		groups := strings.Join(t.Groups, ", ")
		row := []string{
			name,
			t.Language,
			t.Type,
			t.Model,
			groups,
			t.Short,
		}

		rows = append(rows, row)
	}

	return rows
}

func (t *templatePrinter) BuildHeader() []string {

	row := []string{
		"Name",
		"Language",
		"Type",
		"Model",
		"Groups",
		"Description",
	}

	return row
}

type templateGrouper interface {
	group(t map[string]tpl.Template) map[string]map[string]tpl.Template
}

func getGrouper(t string) templateGrouper {

	switch t {

	case "language":
		return &groupByLanguage{}
	case "group":
		return &groupByGroup{}
	case "key":
		return &groupByKey{}
	case "model":
		return &groupByModel{}
	case "tag":
		return &groupByTag{}
	default:
		return &groupByType{}

	}
}

type groupByType struct{}
type groupByGroup struct{}
type groupByTag struct{}
type groupByLanguage struct{}
type groupByKey struct{}
type groupByModel struct{}

func (g *groupByType) group(t map[string]tpl.Template) map[string]map[string]tpl.Template {
	m := make(map[string]map[string]tpl.Template)

	for tn, template := range t {

		k, f := m[template.Type]

		if !f {
			k = make(map[string]tpl.Template)
		}

		k[tn] = template
		m[template.Type] = k

	}

	return m
}
func (g *groupByKey) group(t map[string]tpl.Template) map[string]map[string]tpl.Template {
	m := make(map[string]map[string]tpl.Template)
	empty := make(map[string]tpl.Template)

	for tn, template := range t {

		if len(template.Key) == 0 {
			empty[tn] = template
		} else {

			k, f := m[template.Key]

			if !f {
				k = make(map[string]tpl.Template)
			}

			k[tn] = template
			m[template.Key] = k
		}

	}

	if len(empty) > 0 {
		m["nokey"] = empty
	}

	return m
}
func (g *groupByModel) group(t map[string]tpl.Template) map[string]map[string]tpl.Template {
	m := make(map[string]map[string]tpl.Template)
	empty := make(map[string]tpl.Template)

	for tn, template := range t {

		if len(template.Model) == 0 {
			empty[tn] = template
		} else {

			k, f := m[template.Model]

			if !f {
				k = make(map[string]tpl.Template)
			}

			k[tn] = template
			m[template.Model] = k
		}

	}

	if len(empty) > 0 {
		m["empty"] = empty
	}

	return m
}

func (g *groupByLanguage) group(t map[string]tpl.Template) map[string]map[string]tpl.Template {
	m := make(map[string]map[string]tpl.Template)
	empty := make(map[string]tpl.Template)

	for tn, template := range t {

		if len(template.Language) > 0 {
			empty[tn] = template
		} else {

			if template.Type != "block" {

				k, f := m[template.Language]

				if !f {
					k = make(map[string]tpl.Template)
				}

				k[tn] = template
				m[template.Language] = k
			}
		}
	}

	if len(empty) > 0 {
		m["empty"] = empty
	}

	return m
}
func (g *groupByTag) group(t map[string]tpl.Template) map[string]map[string]tpl.Template {
	m := make(map[string]map[string]tpl.Template)
	empty := make(map[string]tpl.Template)

	for tn, template := range t {

		if template.Type != "block" {

			var tmpl = template

			if len(template.Tags) == 0 {
				empty[tn] = tmpl
			} else {
				for _, grp := range template.Tags {

					k, f := m[grp]

					if !f {
						k = make(map[string]tpl.Template)
					}

					k[tn] = tmpl
					m[grp] = k
				}
			}
		}

	}

	if len(empty) > 0 {
		m["empty"] = empty
	}
	return m
}
func (g *groupByGroup) group(t map[string]tpl.Template) map[string]map[string]tpl.Template {
	m := make(map[string]map[string]tpl.Template)
	empty := make(map[string]tpl.Template)

	for tn, template := range t {

		if template.Type != "block" {

			var tmpl = template

			if len(template.Groups) == 0 {
				empty[tn] = tmpl
			} else {
				for _, grp := range template.Groups {

					k, f := m[grp]

					if !f {
						k = make(map[string]tpl.Template)
					}

					k[tn] = tmpl
					m[grp] = k
				}
			}
		}
	}

	if len(empty) > 0 {
		m["empty"] = empty
	}

	return m
}

type templateFilter interface {
	filter(t map[string]tpl.Template, filter []string) map[string]tpl.Template
}
type filterByLang struct{}
type filterByType struct{}
type filterByModel struct{}
type filterByKey struct{}

func (f *filterByType) filter(t map[string]tpl.Template, filter []string) map[string]tpl.Template {
	output := make(map[string]tpl.Template)

	for name, template := range t {
		if contains(filter, template.Type) {
			output[name] = template
		}
	}

	return output
}
func (f *filterByLang) filter(t map[string]tpl.Template, filter []string) map[string]tpl.Template {
	output := make(map[string]tpl.Template)

	for name, template := range t {
		if contains(filter, template.Language) {
			output[name] = template
		}
	}

	return output
}
func (f *filterByModel) filter(t map[string]tpl.Template, filter []string) map[string]tpl.Template {
	output := make(map[string]tpl.Template)

	for name, template := range t {
		if contains(filter, template.Model) {
			output[name] = template
		}
	}

	return output
}
func (f *filterByKey) filter(t map[string]tpl.Template, filter []string) map[string]tpl.Template {
	output := make(map[string]tpl.Template)

	for name, template := range t {
		if contains(filter, template.Key) {
			output[name] = template
		}
	}

	return output
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
