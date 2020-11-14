package data

import (
	"database/sql"
	"modelhelper/cli/types"
)

type Database interface {
	Connect(source string) (*sql.DB, error)

	Entity(name string) (*types.Entity, error)
	Entities(pattern string) (*[]types.Entity, error)

	Columns(entityName string) (*[]types.EntityColumn, error)

	ParentRelations(entityName string) (*[]types.EntityRelation, error)
	ChildRelations(entityName string) (*[]types.EntityRelation, error)

	Indexes(entityName string) (*[]types.EntityIndex, error)
}

type Optimize interface {
}
