package cmd

import (
	"fmt"
	"sort"
	"strings"

	"modelhelper/cli/app"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/slice"
	"modelhelper/cli/source"
	"modelhelper/cli/tree"
	"modelhelper/cli/ui"

	"github.com/gookit/color"
	_ "github.com/gookit/color"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var skipDescription bool

type entityHeader []string
type entitiesTableRenderer struct {
	rows     []modelhelper.Entity
	withDesc bool
	withStat bool
}

type entityFilter interface {
	Filter(t []modelhelper.Entity, filter []string) []modelhelper.Entity
}

var (
	nodeTable = map[int]*Node{}
	root      *Node
)

// entityCmd represents the entity command
var entityCmd = &cobra.Command{
	Use:     "entity",
	Aliases: []string{"e"},
	Short:   "Show a list of entities or details of a single entity",

	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("entity called")
		var con modelhelper.Connection
		var conName string

		isDemo, _ := cmd.Flags().GetBool("demo")
		conNameFlag, _ := cmd.Flags().GetString("connection")

		modelHelperApp = app.New()
		ctx := modelHelperApp.CreateContext()

		if isDemo {
			conName = "demo"
			con = modelhelper.Connection{Type: conName}
		} else {

			if len(ctx.Connections) == 0 {
				fmt.Println("Could not find any connections to use, please add a connection")
				fmt.Println("to the config and/or any project file")
				return
			}

			if len(conNameFlag) == 0 {

				conName = ctx.DefaultConnection
			} else {
				conName = conNameFlag
			}

			if len(conName) == 0 {
				ka := keyArray(ctx.Connections)
				conName = ka[0]
			}

			con = ctx.Connections[conName]

		}

		src := source.SourceFactory(&con)

		pattern := ""

		if src == nil {
			fmt.Println("Could not load the source, check configuration")
			return
		}
		isSearch := false

		if len(args) > 0 {
			isSearch = isSearchPattern(args[0])

			if isSearch {
				pattern = args[0]
			}
		}
		if len(args) > 0 && !isSearch {
			en := args[0]
			e, err := src.Entity(en)
			if err != nil {
				fmt.Println(err)
			}

			if e == nil {
				fmt.Println("The entity could not be found")
				return
			}

			p := message.NewPrinter(language.English)

			// maxL := len(e.Schema) + len(e.Name)

			fmt.Printf("\nEntity:         %s.%s", e.Schema, e.Name)
			fmt.Printf("\nRows:           %v", p.Sprintf("%d", e.RowCount))
			fmt.Printf("\nIs Versioned:   %s", yesNo(e.IsVersioned))

			if e.IsVersioned {
				fmt.Printf("\nHist. Table:    %s", e.HistoryTable)
			}

			// fmt.Printf("\nCreated: %s\n", "Unknown")

			if len(e.Description) > 0 {
				ui.PrintConsoleTitle("Description:")
				fmt.Println(e.Description)
			}
			renderColumns(&e.Columns)

			if len(e.Indexes) > 0 {
				ui.PrintConsoleTitle("Indexes")

				itr := indexTableRenderer{
					rows: e.Indexes,
				}

				ui.RenderTable(&itr)
			}

			if len(e.ChildRelations) > 0 {
				ui.PrintConsoleTitle("One to many (.ChildRelations)")
				crtr := relTableRenderer{
					rows: e.ChildRelations,
				}

				ui.RenderTable(&crtr)
			}

			if len(e.ParentRelations) > 0 {
				ui.PrintConsoleTitle("Many to one (.ParentRelations)")
				crtr := relTableRenderer{
					rows: e.ParentRelations,
				}
				ui.RenderTable(&crtr)

			}
			fmt.Println("")

			showTree, _ := cmd.Flags().GetBool("tree")
			if showTree {
				// treeLoader := con.LoadRelationTree()

				// if treeLoader != nil {

				// 	flat, err := treeLoader.GetParentRelationTree(e.Schema, e.Name)
				// 	// flatChild, err := treeLoader.GetChildRelationTree(e.Schema, e.Name)
				// 	if err != nil {

				// 	}

				// 	tb := source.RelationTreeBuilder{
				// 		Items: *flat,
				// 	}

				// 	tree.Print(&tb, true)

				// 	// for _, node := range *flat {
				// 	// 	add(node.ID, node.ParentID, node.TableName, node.ColumnName, node.RelatedTable, node.RelatedColumnName)
				// 	// }
				// 	// // for _, node := range *flatChild {
				// 	// // 	add(node.ID, node.ParentID, node.TableName, node.ColumnName, node.RelatedTable, node.RelatedColumnName)
				// 	// // }

				// 	// show()
				// }

			}
		} else {

			columnFilter, _ := cmd.Flags().GetString("column")
			var ents *[]modelhelper.Entity
			var err error

			if len(columnFilter) > 0 {
				ents, err = src.EntitiesFromColumn(columnFilter)
			} else {

				ents, err = src.Entities(pattern)
			}

			if err != nil {
				fmt.Println(err)
			}
			if ents == nil {
				return
			}

			includeDesc, _ := cmd.Flags().GetBool("desc")

			etr := entitiesTableRenderer{
				rows:     *ents,
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

	},
}

type relationTreeBuilder struct {
	items []source.RelationTreeItem
}

func keyArray(input map[string]modelhelper.Connection) []string {
	keys := []string{}
	for k := range input {
		keys = append(keys, k)
	}

	return keys
}

func yesNo(eval bool) string {
	if eval {
		return "Yes"
	}

	return "No"
}

// func (cmd *cobra.Command) printEntities(pattern string) {

// }

func isSearchPattern(input string) bool {
	return strings.ContainsAny(input, "*%")
}
func init() {
	rootCmd.AddCommand(entityCmd)
	entityCmd.Flags().Bool("demo", false, "Uses the demo source")

	entityCmd.Flags().BoolVarP(&skipDescription, "skip-description", "", false, "Does not show description")
	// entityCmd.Flags().String("by", "", "Groups the list of entities by type (view, table), schema")
	entityCmd.Flags().StringArray("type", []string{}, "Filter the entities by the name of the type [view, table]")
	entityCmd.Flags().StringArray("schema", []string{}, "Filter the templates by the name of the schema [dbo, nn]")
	entityCmd.Flags().Bool("desc", false, "Show or hide description (default true)")
	entityCmd.Flags().Bool("has-rows", false, "Filter only entities with rows")
	entityCmd.Flags().Bool("with-relations", false, "Filter only entities with relations")
	entityCmd.Flags().Bool("no-relations", false, "Filter only entities without relations")
	entityCmd.Flags().Bool("is-versioned", false, "Filter only entities that is versioned")

	// entityCmd.Flags().Bool("tree", false, "Filter only entities that is versioned")

	entityCmd.Flags().String("key", "", "The key to use when encoding and decoding secrets for a connection")
	entityCmd.Flags().String("column", "", "List entities based on a column name")
	entityCmd.Flags().StringP("connection", "c", "", "The connection to be used, uses default connection if not provided")

	entityCmd.Flags().String("sort", "name", "Sorts the table values [name, rows], default value: name")
	entityCmd.Flags().Bool("descending", false, "Sorts the list of entities descending")
	// entityCmd.Flags().Bool("include-history", false, "Includes history enities in the list")

}

func renderColumns(cl *modelhelper.ColumnList) {

	ui.PrintConsoleTitle("Columns")

	colr := source.ColumnToTableRenderer{
		IncludeDescription: !skipDescription,
		Columns:            cl,
	}

	ui.RenderTable(&colr)
}

type filterEntityByType struct{}
type filterEntityBySchema struct{}
type filterVersionedEntity struct{}
type filterEntitiesWithRows struct{}
type filterEntitiesWithRelations struct{}

func (f *filterEntitiesWithRows) filter(e []modelhelper.Entity, filter []string) []modelhelper.Entity {
	output := []modelhelper.Entity{}
	for _, entity := range e {
		if entity.RowCount > 0 {
			output = append(output, entity)
		}
	}

	return output
}
func (f *filterEntitiesWithRelations) filter(e []modelhelper.Entity, filter []string) []modelhelper.Entity {
	output := []modelhelper.Entity{}
	for _, entity := range e {
		if (entity.ParentRelationCount + entity.ChildRelationCount) > 0 {
			output = append(output, entity)
		}
	}

	return output
}
func (f *filterVersionedEntity) filter(e []modelhelper.Entity, filter []string) []modelhelper.Entity {
	output := []modelhelper.Entity{}
	for _, entity := range e {
		if entity.IsVersioned {
			output = append(output, entity)
		}
	}

	return output
}
func (f *filterEntityByType) filter(e []modelhelper.Entity, filter []string) []modelhelper.Entity {
	output := []modelhelper.Entity{}
	for _, entity := range e {
		if slice.Contains(filter, entity.Type) {
			output = append(output, entity)
		}
	}

	return output
}
func (f *filterEntityBySchema) filter(e []modelhelper.Entity, filter []string) []modelhelper.Entity {
	output := []modelhelper.Entity{}
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
	rows []modelhelper.Index
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
	rows []modelhelper.Relation
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

func generateRelTree(src *source.Source) map[int]Node {
	return nil
}

type Node struct {
	ID        int
	Name      string
	Column    string
	RelColumn string
	RelName   string
	Nodes     map[int]Node
}

func (tb *relationTreeBuilder) Build() tree.Node {
	root := tree.Node{}
	nodeTable := map[int]tree.Node{}

	add := func(id, parentId int, name, column, relName, relCol string) {
		// internalTable := map[int]tree.Node{}
		desc := color.FgDarkGray.Sprintf("connection: (%s.%s => %s.%s)", name, column, relName, relCol)
		node := tree.Node{Name: name, Description: desc}

		if parentId == -1 {
			root = node
		} else {
			parent, ok := nodeTable[parentId]
			if !ok {
				return
			}

			parent.Nodes = append(parent.Nodes, node)
			//parent.Nodes[id] = node
			// parent.Add(node)
		}

		nodeTable[id] = node
	}

	for _, item := range tb.items {
		add(item.ID, item.ParentID, item.TableName, item.ColumnName, item.RelatedTable, item.RelatedColumnName)
	}

	return root
}

// func (tb *relationTreeBuilder) Describe() string {
// 	desc := color.FgDarkGray.Sprintf("connection: (%s.%s => %s.%s)", name, column, relName, relCol)

// }
// func add(id, parentId int, name, column, relName, relCol string) {
// 	// fmt.Printf("add: id=%v name=%v parentId=%v\n", id, name, parentId)

// 	node := &Node{Name: name, ID: id, Column: column, RelName: relName, RelColumn: relCol, Nodes: make(map[int]Node)}

// 	if parentId == -1 {
// 		root = node
// 	} else {

// 		parent, ok := nodeTable[parentId]
// 		if !ok {
// 			// fmt.Printf("add: parentId=%v: not found\n", parentId)
// 			return
// 		}

// 		parent.Nodes[id] = *node
// 	}

// 	nodeTable[id] = node
// }

func show() {

	ui.PrintConsoleTitle("Relation Tree")
	showNode(*root, "")
}

func showNode(node Node, prefix string) {
	if prefix == "" {
		fmt.Printf("%v\n", node.Name)
	} else {
		// fmt.Printf("%v %v \n", prefix, node.Name)
		// fmt.Printf("%v %v (%s => %s)\n", prefix, node.Name, node.Column, node.RelColumn)
		cols := color.FgDarkGray.Sprintf("connection: (%s.%s => %s.%s)", node.Name, node.Column, node.RelName, node.RelColumn)
		fmt.Printf("%v %v %s\n", prefix, node.Name, cols)
	}
	for _, n := range node.Nodes {
		showNode(n, prefix+"  ")
	}
}
