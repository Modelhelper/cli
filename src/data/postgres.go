package data

import (
	"database/sql"
	"modelhelper/cli/types"
)

type Postgres struct{}

func (server *Postgres) Connect(source string) (*sql.DB, error) {
	return nil, nil
}

func (server *Postgres) Entity(name string) (*types.Entity, error) {
	return nil, nil
}
func (server *Postgres) Entities(pattern string) (*[]types.Entity, error) {
	return nil, nil
}
