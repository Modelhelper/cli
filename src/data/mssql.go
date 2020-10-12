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
