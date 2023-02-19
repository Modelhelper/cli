package source

import (
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/tree"

	"github.com/gookit/color"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type RelationTreeBuilder struct {
	Items []RelationTreeItem
}

// type SortTableById []Entity
type SortTableByName []modelhelper.Entity
type SortTableByNameDesc []modelhelper.Entity
type SortTableByRows []modelhelper.Entity
type SortTableByRowsDesc []modelhelper.Entity

func (a SortTableByName) Len() int           { return len(a) }
func (a SortTableByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a SortTableByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a SortTableByNameDesc) Len() int           { return len(a) }
func (a SortTableByNameDesc) Less(i, j int) bool { return a[i].Name > a[j].Name }
func (a SortTableByNameDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a SortTableByRows) Len() int           { return len(a) }
func (a SortTableByRows) Less(i, j int) bool { return a[i].RowCount < a[j].RowCount }
func (a SortTableByRows) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a SortTableByRowsDesc) Len() int           { return len(a) }
func (a SortTableByRowsDesc) Less(i, j int) bool { return a[i].RowCount > a[j].RowCount }
func (a SortTableByRowsDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type DefaultEntitiesTableRenderer []modelhelper.Entity
type DescriptiveEntitiesRenderer []modelhelper.Entity
type SimpleEntitiesRenderer []modelhelper.Entity

func (d *DefaultEntitiesTableRenderer) Rows() [][]string {
	return toRows(*d, false, true)
}

func (d *DefaultEntitiesTableRenderer) Header() []string {
	return buildHeader(false, true)
}
func (d *DescriptiveEntitiesRenderer) Rows() [][]string {
	return toRows(*d, true, false)
}

func (d *DescriptiveEntitiesRenderer) Header() []string {
	return buildHeader(true, false)
}
func (d *SimpleEntitiesRenderer) Rows() [][]string {
	return toRows(*d, false, false)
}

func (d *SimpleEntitiesRenderer) Header() []string {
	return buildHeader(false, false)
}

func toRows(input []modelhelper.Entity, withDesc, withStat bool) [][]string {
	rows := [][]string{}
	for _, e := range input {
		p := message.NewPrinter(language.English)

		r := []string{
			e.Name,
			e.Schema,
		}
		if !withDesc {
			r = append(r, e.Type)
			r = append(r, e.Alias)
			r = append(r, p.Sprintf("%d", e.RowCount))
		}
		if withStat {
			r = append(r, p.Sprintf("%d", e.ColumnCount))
			r = append(r, p.Sprintf("%d", e.ParentRelationCount))
			r = append(r, p.Sprintf("%d", e.ChildRelationCount))
		}

		if withDesc {
			r = append(r, e.Description)
		}

		rows = append(rows, r)
	}

	return rows

}

func buildHeader(withDesc, withStat bool) []string {
	h := []string{"Name", "Schema"}

	if !withDesc {
		h = append(h, "Type")
		h = append(h, "Alias")
		h = append(h, "Rows")
	}

	if withStat {
		h = append(h, "Col Cnt")
		h = append(h, "P Relations")
		h = append(h, "C Relations")

	}

	if withDesc {
		h = append(h, "Description")
	}

	return h
}

// var nodeTable map[int]tree.Node

func (tb *RelationTreeBuilder) Build() tree.Node {
	var root tree.Node

	for _, item := range tb.Items {
		if item.ParentID == -1 {
			root = tree.Node{
				ID:          item.ID,
				Name:        item.TableName,
				Description: "",
				Nodes:       []tree.Node{},
			}

			addChildNodes(&root, tb.Items)
		} else {
			break
		}
	}

	return root
}

func addChildNodes(n *tree.Node, nodes []RelationTreeItem) {
	for _, item := range nodes {
		if n.ID == item.ParentID {
			desc := color.FgDarkGray.Sprintf("connection: (%s.%s => %s.%s)", item.TableName, item.ColumnName, item.RelatedTable, item.RelatedColumnName)
			child := tree.Node{
				ID:          item.ID,
				Name:        item.TableName,
				Description: desc,
				Nodes:       []tree.Node{},
			}
			addChildNodes(&child, nodes)
			n.Nodes = append(n.Nodes, child)

		}
	}
}
