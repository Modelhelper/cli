package source

import "modelhelper/cli/modelhelper"

type Postgres struct {
	Connection modelhelper.Connection
}

func (server *Postgres) Entity(name string) (*modelhelper.Entity, error) {
	return nil, nil
}
func (server *Postgres) Entities(pattern string) (*[]modelhelper.Entity, error) {
	return nil, nil
}

func (server *Postgres) EntitiesFromColumn(column string) (*[]modelhelper.Entity, error) {
	return nil, nil
}
func (server *Postgres) ConnectionStringPart(part string) string {
	return ""
}
func (server *Postgres) ParseConnectionString() {
}
