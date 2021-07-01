package source

import "testing"

func TestTreeBuilder(t *testing.T) {
	builder := RelationTreeBuilder{
		Items: mockRelations(),
	}

	tree := builder.Build()
	a := len(tree.Nodes)
	e := 2

	if a != e {
		t.Errorf("TreeBuilder: expected %v, got %v", e, a)
	}
}

func mockRelations() []RelationTreeItem {
	data := []RelationTreeItem{
		{KeyName: "Root", ParentID: -1, ID: 1, RelatedTable: "rel", RelatedColumnName: "rel-col", TableName: "tname", ColumnName: "colname"},
		{KeyName: "child 1", ParentID: 1, ID: 2, RelatedTable: "rel", RelatedColumnName: "rel-col", TableName: "tname", ColumnName: "colname"},
		{KeyName: "child 2", ParentID: 1, ID: 3, RelatedTable: "rel", RelatedColumnName: "rel-col", TableName: "tname", ColumnName: "colname"},
		{KeyName: "child 1.1", ParentID: 2, ID: 4, RelatedTable: "rel", RelatedColumnName: "rel-col", TableName: "tname", ColumnName: "colname"},
	}

	return data
}
