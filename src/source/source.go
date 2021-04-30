package source

import "log"

type LanguageDef struct {
	Definitions string
}

type ConnectionMap map[string]Connection
type ConnectionProvider interface {
	GetConnections() (*map[string]Connection, error)
}

type Source interface {
	Entity(name string) (*Entity, error)
	Entities(pattern string) (*[]Entity, error)
}

// should be renamed
type Connection struct {
	Name             string
	Description      string
	ConnectionString string
	Schema           string
	Type             string
	Groups           map[string]ConnectionGroup
	Options          map[string]interface{}
}

// should be renamed
// should this be in the input source package, since it's shared among project, config and other input sources
type ConnectionGroup struct {
	Items   []string
	Options map[string]interface{}
}

type DatabaseOptimizer interface {
	RebuildIndexes()
}

func MergeConnections(providers ...ConnectionProvider) (*map[string]Connection, error) {
	var output = make(map[string]Connection)

	for _, pv := range providers {

		cons, err := pv.GetConnections()
		if err != nil {
			log.Fatal("Could not get connections", err)
		}

		for p, v := range *cons {
			output[p] = v
		}
	}
	return &output, nil
}

func (c *Connection) LoadSource() Source {

	var src Source
	switch c.Type {
	case "mssql":
		src = &MsSql{Connection: *c}
	case "postgres":
		src = &Postgres{Connection: *c}

	case "demo":
		src = &DemoSource{}
	default:
		src = nil
	}

	return src
}

func LoadSource(name string, connections map[string]Connection) Source {

	if name == "demo" {
		return &DemoSource{}
	}

	s, exists := connections[name]

	if !exists {
		return nil
	}

	return s.LoadSource()

}
