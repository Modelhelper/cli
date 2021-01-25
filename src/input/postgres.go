package input

import "modelhelper/cli/app"

type Postgres struct {
	Source app.ConfigSource
}

// func (server *Postgres) Connect(source string) (*sql.DB, error) {
// 	return nil, nil
// }

func (server *Postgres) Entity(name string) (*Entity, error) {
	return nil, nil
}
func (server *Postgres) Entities(pattern string) (*[]Entity, error) {
	return nil, nil
}
