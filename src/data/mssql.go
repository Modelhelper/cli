package data

import (
	"database/sql"
	"fmt"
	"modelhelper/cli/types"
)

// type Database interface {
// 	GetEntity(entityName string)
// }

func GetEntity(entityName string, db *sql.DB) (*types.Entity, error) {
	d := db.Driver()

	fmt.Println(d)
	return nil, nil
}

type MsSql struct{}

func (server *MsSql) Connect(source string) (*sql.DB, error) {
	return nil, nil
}

func (server *MsSql) Entity(name string) (*types.Entity, error) {
	return nil, nil
}
func (server *MsSql) Entities(pattern string) (*[]types.Entity, error) {
	return nil, nil
}
