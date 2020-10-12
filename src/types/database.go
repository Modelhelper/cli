package types

type Entity struct {
	// Creator        Creator
	Name string
	// ModelName      string
	// ContextualName string
	Schema      string
	Type        string
	RowCount    int
	Created     string
	Alias       string
	Description string
	Columns     []EntityColumn

	// Parents     []EntityRelation
	// Children    []EntityRelation
	Indexes []EntityIndex
}

type EntityColumn struct {
	Name string
}

type Database struct {
	Entities   []Entity
	Name       string
	TableCount int
	ViewCount  int
	Server     string
}

type EntityRelation struct {
}
type EntityIndex struct {
}

type IndexColumn struct {
}
