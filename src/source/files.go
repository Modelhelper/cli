package source

import (
	"io/fs"
)

type Files struct {
	Connection Connection
	Path       string
	FileType   string
	FileSystem fs.ReadDirFS
}

type fileEntity struct {
	Name        string                `yaml:"name"`
	Schema      string                `yaml:"schema"`
	Description string                `yaml:"description"`
	Columns     map[string]fileColumn `yaml:"columns"`
}

type fileColumn struct {
	Name        string `yaml:"name"`
	Datatype    string `yaml:"type"`
	Nullable    bool   `yaml:"nullable"`
	References  string `yaml:"references"`
	Identity    bool   `yaml:"identity"`
	Description string `yaml:"description"`
}

/*
Entity(name string) (*Entity, error)
	Entities(pattern string) (*[]Entity, error)
*/
func (f *Files) Entity(name string) (*Entity, error) {
	entities, err := f.Entities("")
	if err != nil {
		return nil, err
	}

	for _, entity := range *entities {
		if entity.Name == name {
			return &entity, nil
		}
	}

	return nil, &EntityNotFoundError{name}
}
func (f *Files) Entities(pattern string) (*[]Entity, error) {
	return nil, nil
}

func (f *fileEntity) toSourceEntity() Entity {
	return Entity{
		Name:        f.Name,
		Schema:      f.Schema,
		Description: f.Description,
		Type:        "table",
		RowCount:    0,
	}
}

// func (f *fileEntity) Open(blob []byte) *fileEntity {

// }
