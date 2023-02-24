package source

import "modelhelper/cli/modelhelper/models"

type Postgres struct {
	Connection models.Connection
}

func (server *Postgres) Entity(name string) (*models.Entity, error) {
	return nil, nil
}
func (server *Postgres) Entities(pattern string) (*[]models.Entity, error) {
	return nil, nil
}
func (server *Postgres) EntitiesFromNames(names []string) (*[]models.Entity, error) {
	return nil, nil
}

func (server *Postgres) EntitiesFromColumn(column string) (*[]models.Entity, error) {
	return nil, nil
}
func (server *Postgres) ConnectionStringPart(part string) string {
	return ""
}
func (server *Postgres) ParseConnectionString() {
}
