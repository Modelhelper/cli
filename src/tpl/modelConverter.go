package tpl

import (
	"modelhelper/cli/source"
	"strings"
)

type EntityToModel struct {
	Entity *source.Entity
}

func (e *EntityToModel) Convert() interface{} {
	model := EntityImportModel{
		Name:              e.Entity.Name,
		Schema:            e.Entity.Schema,
		Parents:           []EntityRelation{},
		Children:          []EntityRelation{},
		Columns:           []EntityColumnImportModel{},
		IgnoredColumns:    []EntityColumnImportModel{},
		NonIgnoredColumns: []EntityColumnImportModel{},
		PrimaryKeys:       []EntityColumnImportModel{},
		ForeignKeys:       []EntityColumnImportModel{},
	}

	for _, c := range e.Entity.Columns {

		col := EntityColumnImportModel{
			Name:           c.Name,
			Description:    c.Description,
			DataType:       c.DataType,
			IsForeignKey:   c.IsForeignKey,
			IsPrimaryKey:   c.IsPrimaryKey,
			IsIdentity:     c.IsIdentity,
			IsNullable:     c.IsNullable,
			HasDescription: len(c.Description) > 0,
		}
		model.Columns = append(model.Columns, col)
		if c.IsPrimaryKey {
			model.PrimaryKeys = append(model.PrimaryKeys, col)
		}

		if c.IsForeignKey {
			model.ForeignKeys = append(model.ForeignKeys, col)
		}

		col.HasPrefix = strings.HasPrefix(c.Name, e.Entity.Name)
		if col.HasPrefix {
			col.NameWithoutPrefix = strings.TrimPrefix(c.Name, e.Entity.Name)
		}

	}

	if len(e.Entity.ChildRelations) > 0 {
		for _, rel := range e.Entity.ChildRelations {
			c := EntityRelation{
				Name: "",
				OwnerColumn: EntityColumnProps{
					Name:       rel.OwnerColumnName,
					DataType:   rel.OwnerColumnType,
					IsNullable: false,
				},
			}

			model.Children = append(model.Children, c)
		}
	}
	return model
}

type EntitiesToModel struct {
	Entities []*EntityImportModel
}

type ProjectToModel struct {
	Project interface{}
}

// func ModelConverter(in interface{}) *interface{} {
// 	t := reflect.TypeOf(in)

// 	switch t.Name() {
// 	case reflect.TypeOf(types.EntityImportModel{}).Name():
// 		e := EntityToModel{}
// 		m, _ := e.Convert()
// 		return m

// 	}
// 	return nil
// }
