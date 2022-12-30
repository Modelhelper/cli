package source

import (
	"io/fs"
	"modelhelper/cli/modelhelper"
	"sort"
)

type Files struct {
	Connection modelhelper.Connection
	Path       string
	FileType   string
	FileSystem fs.ReadDirFS
}

type fileEntity struct {
	Name        string                `yaml:"name"`
	Schema      string                `yaml:"schema"`
	Description string                `yaml:"description"`
	Rows        int                   `yaml:"rows"`
	Columns     map[string]fileColumn `yaml:"columns"`
}

type fileColumn struct {
	ID          int              `yaml:"id"`
	Name        string           `yaml:"name"`
	Datatype    string           `yaml:"type"`
	Nullable    bool             `yaml:"nullable"`
	References  *columnReference `yaml:"references"`
	Identity    bool             `yaml:"identity"`
	IsPrimary   bool             `yaml:"primary"`
	Description string           `yaml:"description"`
}

type columnReference struct {
	Table  string `yaml:"table"`
	Column string `yaml:"column"`
}

type fileEntityMap map[string]fileEntity
type fileEntityList []fileEntity

func (fel *fileEntityList) toMap() map[string]fileEntity {
	m := make(map[string]fileEntity)

	for _, entity := range *fel {
		m[entity.Name] = entity
	}

	return m
}

func (fel *fileEntity) getParentRelations(entities []fileEntity) []modelhelper.Relation {
	relations := []modelhelper.Relation{}

	for colName, col := range fel.Columns {
		if col.References != nil {
			for _, entity := range entities {
				if entity.Name == col.References.Table {
					// col.IsForeignKey = true
					relatedCol, relcolF := entity.Columns[col.References.Column]
					nullable := false
					relcolDT := col.Datatype

					if relcolF {
						nullable = relatedCol.Nullable
						relcolDT = relatedCol.Datatype
					}

					rel := modelhelper.Relation{
						Schema:              entity.Schema,
						ColumnName:          col.References.Column,
						ColumnNullable:      nullable,
						ColumnType:          relcolDT,
						Name:                col.References.Table,
						OwnerColumnName:     colName,
						OwnerColumnNullable: col.Nullable,
						OwnerColumnType:     col.Datatype,
					}

					relations = append(relations, rel)
				}
			}
		}
	}
	return relations
}

func (fel *fileEntity) getChildRelations(entities []fileEntity) []modelhelper.Relation {
	relations := []modelhelper.Relation{}

	for _, entity := range entities {
		if entity.Name != fel.Name {
			for colName, column := range entity.Columns {

				if column.References != nil {
					if column.References.Table == fel.Name {
						// col.IsForeignKey = true
						relatedCol, relcolF := entity.Columns[column.References.Column]
						nullable := false
						relcolDT := column.Datatype

						if relcolF {
							nullable = relatedCol.Nullable
							relcolDT = relatedCol.Datatype
						}

						rel := modelhelper.Relation{
							Schema:              entity.Schema,
							ColumnName:          column.References.Column,
							ColumnNullable:      nullable,
							ColumnType:          relcolDT,
							Name:                entity.Name,
							OwnerColumnName:     colName,
							OwnerColumnNullable: column.Nullable,
							OwnerColumnType:     column.Datatype,
						}

						relations = append(relations, rel)
					}
				}
			}
		}
	}
	return relations
}

func (f *Files) Entity(name string) (*modelhelper.Entity, error) {
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
func (f *Files) Entities(pattern string) (*[]modelhelper.Entity, error) {
	return nil, nil
}
func (f *Files) EntitiesFromColumn(column string) (*[]modelhelper.Entity, error) {
	return nil, nil
}

func (f *fileEntity) toSourceEntity() modelhelper.Entity {
	ent := modelhelper.Entity{
		Name:        f.Name,
		Schema:      f.Schema,
		Description: f.Description,
		Type:        "table",
		RowCount:    f.Rows,
	}

	// parents := []Relation{}
	cols := []modelhelper.Column{}
	id := 0

	for colName, col := range f.Columns {
		id++
		colId := id
		if col.ID > 0 {
			colId = col.ID
		}
		ec := modelhelper.Column{
			ID:           colId,
			Name:         colName,
			DataType:     col.Datatype,
			Description:  col.Description,
			IsNullable:   col.Nullable,
			IsPrimaryKey: col.IsPrimary,
			IsIdentity:   col.Identity,
		}

		if col.References != nil {
			ec.IsForeignKey = true
			// r := Relation{
			// 	Schema:              ent.Schema,
			// 	ColumnName:          col.References.Column,
			// 	Name:                col.References.Table,
			// 	OwnerColumnName:     colName,
			// 	OwnerColumnNullable: col.Nullable,
			// }
			// parents = append(parents, r)
		}

		cols = append(cols, ec)
	}
	sort.Sort(SortColumnById(cols))

	ent.Columns = cols
	return ent
}

func (server *fileEntity) ConnectionStringPart(part string) string {
	return ""
}
func (server *fileEntity) ParseConnectionString() {
}
