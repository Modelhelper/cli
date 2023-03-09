package files

import (
	"fmt"
	"io/fs"
	"modelhelper/cli/modelhelper/models"
	"sort"
)

type Files struct {
	Connection models.Connection
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

func (fel *fileEntity) getParentRelations(entities []fileEntity) []models.Relation {
	relations := []models.Relation{}

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

					rel := models.Relation{
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

func (fel *fileEntity) getChildRelations(entities []fileEntity) []models.Relation {
	relations := []models.Relation{}

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

						rel := models.Relation{
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

func (f *Files) Entity(name string) (*models.Entity, error) {
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

type EntityNotFoundError struct {
	Name string
}

func (e *EntityNotFoundError) Error() string {
	return fmt.Sprintf("Entity '%s' not found", e.Name)
}

func (f *Files) Entities(pattern string) (*[]models.Entity, error) {
	return nil, nil
}
func (f *Files) EntitiesFromColumn(column string) (*[]models.Entity, error) {
	return nil, nil
}

func (f *fileEntity) toSourceEntity() models.Entity {
	ent := models.Entity{
		Name:        f.Name,
		Schema:      f.Schema,
		Description: f.Description,
		Type:        "table",
		RowCount:    f.Rows,
	}

	// parents := []Relation{}
	cols := []models.Column{}
	id := 0

	for colName, col := range f.Columns {
		id++
		colId := id
		if col.ID > 0 {
			colId = col.ID
		}
		ec := models.Column{
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

type SortColumnById []models.Column
type SortColumnByName []models.Column
type ColumnToTableRenderer struct {
	IncludeDescription bool
	Columns            *models.ColumnList
}

func (d *ColumnToTableRenderer) Rows() [][]string {
	var rows [][]string

	for _, c := range *d.Columns {

		null := "No"
		if c.IsNullable {
			null = "Yes"
		}
		id := ""
		if c.IsIdentity {
			id = "Yes"
		}

		pk := ""
		if c.IsPrimaryKey {
			pk = "Yes"
		}

		fk := ""
		if c.IsForeignKey {
			fk = "Yes"
		}

		r := []string{
			c.Name,
			c.DbType,
			null,
			id,
			pk,
			fk,
		}

		if d.IncludeDescription {
			r = append(r, c.Description)
		}
		rows = append(rows, r)
	}

	return rows

}

func (d *ColumnToTableRenderer) Header() []string {
	h := []string{"Name", "Type", "Nullable", "Identity", "PK", "FK"}

	if d.IncludeDescription {
		h = append(h, "Description") //"Children", "Parents"
	}

	return h
}

func (a SortColumnById) Len() int           { return len(a) }
func (a SortColumnById) Less(i, j int) bool { return a[i].ID < a[j].ID }
func (a SortColumnById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a SortColumnByName) Len() int           { return len(a) }
func (a SortColumnByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a SortColumnByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
