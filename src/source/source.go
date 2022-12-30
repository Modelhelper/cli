package source

import (
	"fmt"
	"modelhelper/cli/modelhelper"
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

type ConnectionMap map[string]modelhelper.Connection

// type ConnectionProvider interface {
// 	GetConnections() (*map[string]modelhelper.Connection, error)
// }

type Source interface {
	Entity(name string) (*modelhelper.Entity, error)
	Entities(pattern string) (*[]modelhelper.Entity, error)
	EntitiesFromColumn(column string) (*[]modelhelper.Entity, error)
}

type RelationTree interface {
	GetParentRelationTree(schema string, entityName string) (*[]RelationTreeItem, error)
	GetChildRelationTree(schema string, entityName string) (*[]RelationTreeItem, error)
	// Entities(pattern string) (*[]Entity, error)
}

// should be renamed
// type Connection struct {
// 	Name             string                     `json:"name" yaml:"name"`
// 	Description      string                     `json:"description" yaml:"description,omitempty"`
// 	ConnectionString string                     `json:"connectionString" yaml:"connectionString"`
// 	Schema           string                     `json:"schema" yaml:"schema"`
// 	Database         string                     `json:"database,omitempty" yaml:"database,omitempty"`
// 	Server           string                     `json:"server,omitempty" yaml:"server,omitempty"`
// 	Type             string                     `json:"type" yaml:"type"`
// 	Port             int                        `json:"port,omitempty" yaml:"port,omitempty"`
// 	Entities         []string                   `json:"entities,omitempty" yaml:"entities,omitempty"`
// 	Groups           map[string]ConnectionGroup `json:"groups,omitempty" yaml:"groups,omitempty"`
// 	Options          map[string]interface{}     `json:"options,omitempty" yaml:"options,omitempty"`
// 	Synonyms         map[string]string          `json:"synonyms,omitempty" yaml:"synonyms,omitempty"`
// }

// // should be renamed
// // should this be in the input source package, since it's shared among project, config and other input sources
// type ConnectionGroup struct {
// 	Items   []string               `json:"items" yaml:"items"`
// 	Options map[string]interface{} `json:"options" yaml:"options"`
// }

// type Synonym struct {
// 	Name string
// }

type DatabaseOptimizer interface {
	RebuildIndexes()
}

func BuildConnectionstring(dbtype, server, database, username, password string) string {
	constr := ""
	if dbtype == "mssql" {
		return fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", username, password, server, database)
	}

	return constr
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

// func (c *Connection) ParseConnectionString() {
// 	items := SplitConnectionString(c.ConnectionString)

// 	c.Server = items["server"]
// 	c.Database = items["database"]

// }

// func (c *Connection) ConnectionStringPart(part string) string {
// 	items := SplitConnectionString(c.ConnectionString)

// 	// if out, found := items[part]; f {
// 	if out, found := items[part]; found {
// 		return out
// 	}

// 	return ""
// }

func IsConnectionTypeValid(t string) bool {
	valid := make(map[string]string)
	valid["mssql"] = "Connects to a Microsoft SQL Server"

	_, f := valid[strings.ToLower(t)]
	return f
}

// func MergeConnections(providers ...modelhelper.ConnectionProvider) (*map[string]modelhelper.Connection, error) {
// 	var output = make(map[string]modelhelper.Connection)

// 	for _, pv := range providers {

// 		cons, err := pv.GetConnections()
// 		if err != nil {
// 			log.Fatal("Could not get connections", err)
// 		}

// 		for p, v := range *cons {
// 			current, found := output[p]
// 			if found {
// 				if current.Type == v.Type && len(current.ConnectionString) > 0 && len(v.ConnectionString) == 0 {
// 					v.ConnectionString = current.ConnectionString
// 				}
// 			}

// 			output[p] = v
// 		}
// 	}
// 	return &output, nil
// }

// JoinConnections will merge or replace all the connections it is given
// It works from left to right
// joiner = merge | replace | empty | mergereplace
// func JoinConnections(joinMethod string, connections ...modelhelper.ConnectionProvider) map[string]modelhelper.Connection {
// 	switch joinMethod {
// 	case "merge":
// 		return mergeConnections(connections...)
// 	case "smart":
// 		return smartMergeConnections(connections...)
// 	case "replace":
// 		return replaceConnections(connections...)
// 	default:
// 		return mergeConnections(connections...)
// 	}
// }

// func mergeConnections(connections ...modelhelper.ConnectionProvider) map[string]modelhelper.Connection {
// 	output := make(map[string]modelhelper.Connection)

// 	for _, pv := range connections {
// 		cons, err := pv.GetConnections()
// 		if err != nil {
// 			log.Fatal("Could not get connections", err)
// 		}

// 		for p, v := range *cons {
// 			current, found := output[p]
// 			if found {
// 				if current.Type == v.Type && len(current.ConnectionString) > 0 && len(v.ConnectionString) == 0 {
// 					v.ConnectionString = current.ConnectionString
// 				}
// 			}

// 			output[p] = v
// 		}
// 	}
// 	return output
// }
// func smartMergeConnections(connections ...modelhelper.ConnectionProvider) map[string]modelhelper.Connection {
// 	output := make(map[string]modelhelper.Connection)

// 	for _, pv := range connections {

// 		cons, err := pv.GetConnections()
// 		if err != nil {
// 			log.Fatal("Could not get connections", err)
// 		}

// 		for p, v := range *cons {
// 			current, found := output[p]
// 			if found {
// 				if current.Type == v.Type && len(current.ConnectionString) > 0 && len(v.ConnectionString) == 0 {
// 					v.ConnectionString = current.ConnectionString
// 				}
// 			}

// 			output[p] = v
// 		}
// 	}
// 	return output
// }
// func replaceConnections(connections ...modelhelper.ConnectionProvider) map[string]modelhelper.Connection {
// 	output := make(map[string]modelhelper.Connection)

// 	for _, pv := range connections {

// 		cons, err := pv.GetConnections()
// 		if err != nil {
// 			log.Fatal("Could not get connections", err)
// 		}

// 		for p, v := range *cons {
// 			output[p] = v
// 		}
// 	}
// 	return output
// }

func SourceFactory(c *modelhelper.Connection) modelhelper.SourceService {

	// }

	// func (c *Connection) LoadSource() Source {

	var src modelhelper.SourceService
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

// func (c *Connection) LoadRelationTree() RelationTree {

// 	var src RelationTree
// 	switch c.Type {
// 	case "mssql":
// 		src = &MsSql{Connection: *c}
// 	case "postgres":
// 		src = nil

// 	case "demo":
// 		src = nil
// 	default:
// 		src = nil
// 	}

// 	return src
// }

// func LoadSource(name string, connections map[string]modelhelper.Connection) Source {

// 	if name == "demo" {
// 		return &DemoSource{}
// 	}

// 	s, exists := connections[name]

// 	if !exists {
// 		return nil
// 	}

// 	return s.LoadSource()

// }

func Abbreviate(s string) string {
	abr := ""
	for i, c := range s {
		if i == 0 || unicode.IsUpper(c) {
			abr = abr + string(c)
		}
	}

	return strings.ToLower(abr)
}
