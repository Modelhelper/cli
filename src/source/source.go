package source

import (
	"log"
	"strings"
)

type LanguageDef struct {
	Definitions string
}

type ConnectionMap map[string]Connection
type ConnectionProvider interface {
	GetConnections() (*map[string]Connection, error)
}

type Source interface {
	Entity(name string) (*Entity, error)
	Entities(pattern string) (*EntityList, error)
}

// should be renamed
type Connection struct {
	Name             string                     `json:"name" yaml:"name"`
	Description      string                     `json:"description" yaml:"description"`
	ConnectionString string                     `json:"connectionString" yaml:"connectionString"`
	Schema           string                     `json:"schema" yaml:"schema"`
	Type             string                     `json:"type" yaml:"type"`
	Groups           map[string]ConnectionGroup `json:"groups" yaml:"groups"`
	Options          map[string]interface{}     `json:"options" yaml:"options"`
}

// should be renamed
// should this be in the input source package, since it's shared among project, config and other input sources
type ConnectionGroup struct {
	Items   []string               `yaml:"items" yaml:"items"`
	Options map[string]interface{} `yaml:"options" yaml:"options"`
}

type DatabaseOptimizer interface {
	RebuildIndexes()
}

func IsConnectionTypeValid(t string) bool {
	valid := make(map[string]string)
	valid["mssql"] = "Connects to a Microsoft SQL Server"

	_, f := valid[strings.ToLower(t)]
	return f
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

func (c *Connection) LoadRelationTree() RelationTree {

	var src RelationTree
	switch c.Type {
	case "mssql":
		src = &MsSql{Connection: *c}
	case "postgres":
		src = nil

	case "demo":
		src = nil
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
