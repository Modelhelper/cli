package source

import (
	"fmt"
	"testing"
)

func TestTreeBuilder(t *testing.T) {
	builder := RelationTreeBuilder{
		Items: mockRelations(),
	}

	tree := builder.Build()
	a := len(tree.Nodes)
	e := 2

	for i, node := range tree.Nodes {
		fmt.Printf("\n PARENT:: IDx: %d, name: %s, ln: %d", i, node.Name, len(node.Nodes))

		for i, node := range node.Nodes {
			fmt.Printf("\n CHILD:: IDx: %d, name: %s, ln: %d", i, node.Name, len(node.Nodes))
		}
	}

	fmt.Println(tree)
	if a != e {
	}
	t.Errorf("TreeBuilder: expected %v, got %v", e, a)

	// if len(tree.Nodes[1].Nodes) != 1 {
	// 	t.Errorf("TreeBuilder: expected %v, got %v", 1, len(tree.Nodes[1].Nodes))

	// }
}

func mockRelations() []RelationTreeItem {
	data := []RelationTreeItem{
		{ParentID: -1, ID: 1, RelatedTable: "rel", RelatedColumnName: "rel-col", TableName: "root", ColumnName: "colname"},
		{ParentID: 1, ID: 2, RelatedTable: "rel", RelatedColumnName: "rel-col", TableName: "table 1", ColumnName: "colname"},
		{ParentID: 1, ID: 3, RelatedTable: "rel", RelatedColumnName: "rel-col", TableName: "table 2", ColumnName: "colname"},
		{ParentID: 2, ID: 4, RelatedTable: "rel", RelatedColumnName: "rel-col", TableName: "table 1 - 1", ColumnName: "colname"},
		{ParentID: 3, ID: 5, RelatedTable: "rel", RelatedColumnName: "rel-col", TableName: "table 2 - 1", ColumnName: "colname"},
	}

	return data
}
