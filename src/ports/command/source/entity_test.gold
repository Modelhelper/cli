package cmd_test

import (
	"modelhelper/cli/source"
	"testing"
)

func TestTreeBuilder(t *testing.T) {
	// builder := relationTreeBuilder{
	// 	items: mockRelations(),
	// }

	// tree := builder.Build()

	if len("tree.Nodes") != 2 {
		// t.Errorf("Err %s", "tree")
	}
}

func mockRelations() []source.RelationTreeItem {
	data := []source.RelationTreeItem{
		{KeyName: "Root", ParentID: -1, ID: 1, RelatedTable: "rel", RelatedColumnName: "rel-col", TableName: "tname", ColumnName: "colname"},
		{KeyName: "child 1", ParentID: 1, ID: 2, RelatedTable: "rel", RelatedColumnName: "rel-col", TableName: "tname", ColumnName: "colname"},
		{KeyName: "child 2", ParentID: 1, ID: 3, RelatedTable: "rel", RelatedColumnName: "rel-col", TableName: "tname", ColumnName: "colname"},
		{KeyName: "child 1.1", ParentID: 2, ID: 4, RelatedTable: "rel", RelatedColumnName: "rel-col", TableName: "tname", ColumnName: "colname"},
	}

	return data
}
