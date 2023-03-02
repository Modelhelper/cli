package source

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"modelhelper/cli/ports/source"
	"modelhelper/cli/ui"
	"modelhelper/cli/utils/slice"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func ListCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all source items",

		Run: listTemplateCommandHandler(app),
	}

	addFlags(cmd)

	return cmd
}

func addFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("connection", "c", "", "Groups the templates by type, group, language, model or tag")
	cmd.Flags().String("by", "", "Groups the templates by type, group, language, model or tag")
	cmd.Flags().StringArray("type", []string{}, "Filter the templates by the name of the type")
	cmd.Flags().StringArray("lang", []string{}, "Filter the templates by language")
	cmd.Flags().StringArray("model", []string{}, "Filter the templates by model")
	cmd.Flags().StringArray("key", []string{}, "Filter the templates by key")
	cmd.Flags().StringArray("group", []string{}, "Filter the templates by group")
	cmd.Flags().Bool("demo", false, "Use demo models or not")

	cmd.Flags().String("column", "", "List entities based on a column name")
	cmd.Flags().String("sort", "name", "Sorts the table values [name, rows], default value: name")
	cmd.Flags().Bool("descending", false, "Sorts the list of entities descending")
	cmd.Flags().StringArray("schema", []string{}, "Filter the templates by the name of the schema [dbo, nn]")

}

func listTemplateCommandHandler(app *modelhelper.ModelhelperCli) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		options := parseArgs(cmd, args)
		conName, conType := "", ""
		pattern := ""
		isSearch := false

		if options.IsDemo {
			options.ConnectionName = "demo"
			conName = "demo"
			conType = "file"

		} else {

			connections, err := app.ConnectionService.Connections()
			if err != nil {
				// return nil, err
			}
			if len(connections) == 0 {
				// return nil, errors.New("Could not find any connections to use, please add a connection to the config file")
			}
			if len(options.ConnectionName) == 0 {

				options.ConnectionName = app.Config.DefaultConnection
			}

			if len(options.ConnectionName) == 0 {
				for _, v := range connections {
					options.ConnectionName = v.Name
					break
				}
			}

			conName = options.ConnectionName
			conType = connections[conName].Type
			// con = g.connectionService.Connection(options.ConnectionName)
		}

		src, _ := app.SourceFactory.NewSource(conType, conName)

		if len(args) > 0 {
			isSearch = isSearchPattern(args[0])

			if isSearch {
				pattern = args[0]
			}
		}

		entities, err := src.Entities(pattern)

		if err != nil {
			fmt.Println(err)
		}
		if entities == nil {
			return
		}

		includeDesc, _ := cmd.Flags().GetBool("desc")

		etr := entitiesTableRenderer{
			rows:     *entities,
			withDesc: includeDesc,
			withStat: !includeDesc,
		}
		typeFiler, _ := cmd.Flags().GetStringArray("type")

		if len(typeFiler) > 0 {
			ft := filterEntityByType{}
			etr.rows = ft.filter(etr.rows, typeFiler)
		}

		schemaFilter, _ := cmd.Flags().GetStringArray("schema")
		if len(schemaFilter) > 0 {
			ft := filterEntityBySchema{}
			etr.rows = ft.filter(etr.rows, schemaFilter)
		}

		hasRowsFilter, _ := cmd.Flags().GetBool("has-rows")
		if hasRowsFilter {
			ft := filterEntitiesWithRows{}
			etr.rows = ft.filter(etr.rows, nil)
		}
		hasRelationsFilter, _ := cmd.Flags().GetBool("with-relations")
		if hasRelationsFilter {
			ft := filterEntitiesWithRelations{}
			etr.rows = ft.filter(etr.rows, nil)
		}

		// sorting
		sorting, _ := cmd.Flags().GetString("sort")
		descending, _ := cmd.Flags().GetBool("descending")
		if sorting == "name" {
			if descending {
				sort.Sort(source.SortTableByNameDesc(etr.rows))
			} else {
				sort.Sort(source.SortTableByName(etr.rows))

			}
		} else if sorting == "rows" || sorting == "row" || sorting == "rowcount" {
			if descending {
				sort.Sort(source.SortTableByRowsDesc(etr.rows))
			} else {
				sort.Sort(source.SortTableByRows(etr.rows))

			}
		}

		ui.RenderTable(&etr)

	}
}

func isSearchPattern(input string) bool {
	return strings.ContainsAny(input, "*%")
}

func parseArgs(cmd *cobra.Command, args []string) *models.SourceListOptions {
	connection, _ := cmd.Flags().GetString("connection")
	demo, _ := cmd.Flags().GetBool("demo")
	group, _ := cmd.Flags().GetString("by")
	typeFiler, _ := cmd.Flags().GetStringArray("type")
	langFilter, _ := cmd.Flags().GetStringArray("lang")
	modelFilter, _ := cmd.Flags().GetStringArray("model")
	keyFilter, _ := cmd.Flags().GetStringArray("key")
	groupFilter, _ := cmd.Flags().GetStringArray("group")

	options := &models.SourceListOptions{
		ConnectionName:  connection,
		GroupBy:         group,
		FilterTypes:     typeFiler,
		FilterLanguages: langFilter,
		FilterModels:    modelFilter,
		FilterKeys:      keyFilter,
		FilterGroups:    groupFilter,
		IsDemo:          demo,
	}

	return options
}

type entityHeader []string
type entitiesTableRenderer struct {
	rows     []models.Entity
	withDesc bool
	withStat bool
}

type entityFilter interface {
	Filter(t []models.Entity, filter []string) []models.Entity
}

func renderColumns(cl *models.ColumnList) {

	ui.PrintConsoleTitle("Columns")

	colr := source.ColumnToTableRenderer{
		// IncludeDescription: !skipDescription,
		IncludeDescription: false,
		Columns:            cl,
	}

	ui.RenderTable(&colr)
}

type filterEntityByType struct{}
type filterEntityBySchema struct{}
type filterVersionedEntity struct{}
type filterEntitiesWithRows struct{}
type filterEntitiesWithRelations struct{}

func (f *filterEntitiesWithRows) filter(e []models.Entity, filter []string) []models.Entity {
	output := []models.Entity{}
	for _, entity := range e {
		if entity.RowCount > 0 {
			output = append(output, entity)
		}
	}

	return output
}
func (f *filterEntitiesWithRelations) filter(e []models.Entity, filter []string) []models.Entity {
	output := []models.Entity{}
	for _, entity := range e {
		if (entity.ParentRelationCount + entity.ChildRelationCount) > 0 {
			output = append(output, entity)
		}
	}

	return output
}
func (f *filterVersionedEntity) filter(e []models.Entity, filter []string) []models.Entity {
	output := []models.Entity{}
	for _, entity := range e {
		if entity.IsVersioned {
			output = append(output, entity)
		}
	}

	return output
}
func (f *filterEntityByType) filter(e []models.Entity, filter []string) []models.Entity {
	output := []models.Entity{}
	for _, entity := range e {
		if slice.Contains(filter, entity.Type) {
			output = append(output, entity)
		}
	}

	return output
}
func (f *filterEntityBySchema) filter(e []models.Entity, filter []string) []models.Entity {
	output := []models.Entity{}
	for _, entity := range e {
		if slice.Contains(filter, entity.Schema) {
			output = append(output, entity)
		}
	}

	return output
}

func (d *entitiesTableRenderer) Rows() [][]string {
	var rows [][]string

	for _, e := range d.rows {
		p := message.NewPrinter(language.English)

		r := []string{
			e.Name,
			e.Schema,
		}
		if !d.withDesc {
			r = append(r, e.Type)
			r = append(r, e.Alias)
			r = append(r, p.Sprintf("%d", e.RowCount))
		}
		if d.withStat {
			r = append(r, p.Sprintf("%d", e.ColumnCount))
			r = append(r, p.Sprintf("%d", e.ParentRelationCount))
			r = append(r, p.Sprintf("%d", e.ChildRelationCount))
		}

		if d.withDesc {
			r = append(r, e.Description)
		}

		rows = append(rows, r)
	}

	return rows

}

func (d *entitiesTableRenderer) Header() []string {
	h := []string{"Name", "Schema"}

	if !d.withDesc {
		h = append(h, "Type")
		h = append(h, "Alias")
		h = append(h, "Rows")
	}

	if d.withStat {
		h = append(h, "Col Cnt")
		h = append(h, "P Relations")
		h = append(h, "C Relations")

	}

	if d.withDesc {
		h = append(h, "Description")
	}

	return h
}

type indexTableRenderer struct {
	rows []models.Index
}

func (r *indexTableRenderer) Header() []string {
	return []string{
		"Name",
		"Clustered",
		"Primary",
		"Unique",
		// "Fragmentation",
	}
}

func (r *indexTableRenderer) Rows() [][]string {
	var rows [][]string

	for _, i := range r.rows {
		// p := message.NewPrinter(language.English)

		cluster := "No"
		primary := "No"
		unique := "No"

		if i.IsClustered {
			cluster = "Yes"
		}
		if i.IsUnique {
			unique = "Yes"
		}
		if i.IsPrimaryKey {
			primary = "Yes"
		}
		r := []string{
			i.Name,
			cluster,
			primary,
			unique,
			// p.Sprintf("%d%%", i.AvgFragmentationPercent),
		}

		rows = append(rows, r)
	}

	return rows

}

type relTableRenderer struct {
	rows []models.Relation
}

func (r *relTableRenderer) Header() []string {
	return []string{"Schema", "Name", "ChildCol", "ParentCol", "Constraint"}
}

func (r *relTableRenderer) Rows() [][]string {
	var rows [][]string

	for _, i := range r.rows {
		// p := message.NewPrinter(language.English)

		cn, pn := "NOT NULL", "NOT NULL"

		if i.ColumnNullable {
			cn = "NULL"
		}

		if i.OwnerColumnNullable {
			pn = "NULL"
		}

		r := []string{
			i.Schema,
			i.Name,
			fmt.Sprintf("%s (%s %s)", i.OwnerColumnName, i.OwnerColumnType, pn),
			fmt.Sprintf("%s (%s %s)", i.ColumnName, i.ColumnType, cn),
			i.ContraintName,
		}

		rows = append(rows, r)
	}

	return rows

}
