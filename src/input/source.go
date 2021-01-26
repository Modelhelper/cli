package input

import (
	"modelhelper/cli/app"
)

type Source interface {
	// Connect(source string) (*sql.DB, error)
	// CanConnect() (bool, error)
	Entity(name string) (*Entity, error)
	Entities(pattern string) (*[]Entity, error)

	// Columns(entityName string) (*[]Column, error)

	// ParentRelations(entityName string) (*[]Relation, error)
	// ChildRelations(entityName string) (*[]Relation, error)

	// Indexes(entityName string) (*[]Index, error)
}

type DatabaseOptimizer interface {
	RebuildIndexes()
}

func GetSource(name string, config app.Config) Source {
	s := config.Sources[name]

	if name == "demo" {
		return &DemoInput{}
	}

	var m Source
	switch s.Type {
	case "mssql":
		m = &MsSql{Source: s}
	case "postgres":
		m = &Postgres{Source: s}

	case "demo":
		m = &DemoInput{}
	}

	return m
}
