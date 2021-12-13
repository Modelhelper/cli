package source

import (
	"fmt"
	"log"
	"strings"
	"unicode"
)

type EntityNotFoundError struct {
	Name string
}

func (e *EntityNotFoundError) Error() string {
	return fmt.Sprintf("Entity '%s' not found", e.Name)
}

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
	EntitiesFromColumn(column string) (*[]Entity, error)
}

type RelationTree interface {
	GetParentRelationTree(schema string, entityName string) (*[]RelationTreeItem, error)
	GetChildRelationTree(schema string, entityName string) (*[]RelationTreeItem, error)
	// Entities(pattern string) (*[]Entity, error)
}

// should be renamed
type Connection struct {
	Name             string                     `json:"name" yaml:"name"`
	Description      string                     `json:"description" yaml:"description"`
	ConnectionString string                     `json:"connectionString" yaml:"connectionString"`
	Schema           string                     `json:"schema" yaml:"schema"`
	Database         string                     `json:"database" yaml:"database"`
	Server           string                     `json:"server" yaml:"server"`
	Type             string                     `json:"type" yaml:"type"`
	Port             int                        `json:"port" yaml:"port"`
	Entities         []string                   `json:"entities" yaml:"entities"`
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

func SplitConnectionString(connectionString string) map[string]string {
	items := make(map[string]string)

	if len(connectionString) > 0 {
		parts := strings.Split(connectionString, ";")

		for _, part := range parts {
			kv := strings.Split(part, "=")

			if len(kv) > 0 {
				k, v := strings.ToLower(kv[0]), ""
				if len(kv) == 2 {
					v = kv[1]
				}
				items[k] = v
			}
		}
	}

	return items
}

func (c *Connection) ParseConnectionString() {
	items := SplitConnectionString(c.ConnectionString)

	c.Server = items["server"]
	c.Database = items["database"]

}

func (c *Connection) ConnectionStringPart(part string) string {
	items := SplitConnectionString(c.ConnectionString)

	// if out, found := items[part]; f {
	if out, found := items[part]; found {
		return out
	}

	return ""
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
			current, found := output[p]
			if found {
				if current.Type == v.Type && len(current.ConnectionString) > 0 && len(v.ConnectionString) == 0 {
					v.ConnectionString = current.ConnectionString
				}
			}

			output[p] = v
		}
	}
	return &output, nil
}

//JoinConnections will merge or replace all the connections it is given
//It works from left to right
//joiner = merge | replace | empty | mergereplace
func JoinConnections(joinMethod string, connections ...ConnectionProvider) map[string]Connection {
	switch joinMethod {
	case "merge":
		return mergeConnections(connections...)
	case "smart":
		return smartMergeConnections(connections...)
	case "replace":
		return replaceConnections(connections...)
	default:
		return mergeConnections(connections...)
	}
}

func mergeConnections(connections ...ConnectionProvider) map[string]Connection {
	output := make(map[string]Connection)

	for _, pv := range connections {
		cons, err := pv.GetConnections()
		if err != nil {
			log.Fatal("Could not get connections", err)
		}

		for p, v := range *cons {
			current, found := output[p]
			if found {
				if current.Type == v.Type && len(current.ConnectionString) > 0 && len(v.ConnectionString) == 0 {
					v.ConnectionString = current.ConnectionString
				}
			}

			output[p] = v
		}
	}
	return output
}
func smartMergeConnections(connections ...ConnectionProvider) map[string]Connection {
	output := make(map[string]Connection)

	for _, pv := range connections {

		cons, err := pv.GetConnections()
		if err != nil {
			log.Fatal("Could not get connections", err)
		}

		for p, v := range *cons {
			current, found := output[p]
			if found {
				if current.Type == v.Type && len(current.ConnectionString) > 0 && len(v.ConnectionString) == 0 {
					v.ConnectionString = current.ConnectionString
				}
			}

			output[p] = v
		}
	}
	return output
}
func replaceConnections(connections ...ConnectionProvider) map[string]Connection {
	output := make(map[string]Connection)

	for _, pv := range connections {

		cons, err := pv.GetConnections()
		if err != nil {
			log.Fatal("Could not get connections", err)
		}

		for p, v := range *cons {
			output[p] = v
		}
	}
	return output
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

func Abbreviate(s string) string {
	abr := ""
	for i, c := range s {
		if i == 0 || unicode.IsUpper(c) {
			abr = abr + string(c)
		}
	}

	return strings.ToLower(abr)
}
