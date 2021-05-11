package source

type Postgres struct {
	Connection Connection
}

func (server *Postgres) Entity(name string) (*Entity, error) {
	return nil, nil
}
func (server *Postgres) Entities(pattern string) (*EntityList, error) {
	return nil, nil
}
