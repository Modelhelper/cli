package source

import (
	"modelhelper/cli/tree"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type RelationTreeBuilder struct {
	Items []RelationTreeItem
}

// type SortTableById []Entity
type SortTableByName []Entity
type SortTableByNameDesc []Entity
type SortTableByRows []Entity
type SortTableByRowsDesc []Entity

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

type DefaultEntitiesTableRenderer []Entity
type DescriptiveEntitiesRenderer []Entity
type SimpleEntitiesRenderer []Entity

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

func toRows(input []Entity, withDesc, withStat bool) [][]string {
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
	root := tree.Node{}

	// add := func(id, parentId int, name, column, relName, relCol string) {
	// 	// internalTable := map[int]tree.Node{}
	// 	desc := color.FgDarkGray.Sprintf("connection: (%s.%s => %s.%s)", name, column, relName, relCol)
	// 	node := tree.Node{Name: name, Description: desc}

	// 	if parentId == -1 {
	// 		root = node
	// 	} else {
	// 		parent, ok := nodeTable[parentId]
	// 		if !ok {
	// 			return
	// 		}

	// 		parent.Nodes = append(parent.Nodes, node)
	// 		//parent.Nodes[id] = node
	// 		// parent.Add(node)
	// 	}

	// 	nodeTable[id] = node
	// }

	// for _, item := range tb.Items {
	// 	add(nodeTable, item.ID, item.ParentID, item.TableName, item.ColumnName, item.RelatedTable, item.RelatedColumnName)
	// }

	return root
}

func add(items []RelationTreeItem, id, parentId int, name, column, relName, relCol string) {
	// internalTable := map[int]tree.Node{}
	// desc := color.FgDarkGray.Sprintf("connection: (%s.%s => %s.%s)", name, column, relName, relCol)
	// node := tree.Node{Name: name, Description: desc}

	// if parentId == -1 {
	// 	// root = node
	// } else {
	// 	// parent, ok := nodeTable[parentId]
	// 	// if !ok {
	// 	// 	return
	// 	// }

	// 	parent.Nodes = append(parent.Nodes, node)
	// 	//parent.Nodes[id] = node
	// 	// parent.Add(node)
	// }

	// nodeTable[id] = node
}
