package source

import "embed"

//go:embed entities/*
var entities embed.FS

type DemoSource struct{}

func (server *DemoSource) Entity(name string) (*Entity, error) {

	entities := server.getEntities()
	//var entity Entity

	for _, e := range entities {
		if e.Name == name {
			return &e, nil
		}
	}

	return nil, nil
}
func (server *DemoSource) Entities(pattern string) (*EntityList, error) {
	e := server.getEntities()
	return &e, nil
}

func (server *DemoSource) getEntities() EntityList {
	e := []Entity{

		server.getOrderHeadTable(),
		server.getCustomerTable(),
	}
	return e
}

func (server *DemoSource) getOrderHeadTable() Entity {
	return Entity{
		Name: "order", Schema: "dbo", Description: "This is the order table", RowCount: 1000, Alias: "o", UsesIdentityColumn: true,
		Columns: []Column{
			{Name: "Id", DataType: "int", IsPrimaryKey: true, IsNullable: false, IsIdentity: true, Description: "The identifier"},
			{Name: "Name", DataType: "varchar", IsPrimaryKey: false, IsNullable: false, IsIdentity: false, Description: "Name of the order"},
		},
	}
}
func (server *DemoSource) getCustomerTable() Entity {
	var e = Entity{
		Name: "customer", Schema: "dbo",
		Description: "This is the customer table", RowCount: 1000, Alias: "c", UsesIdentityColumn: true,
		Columns: []Column{
			{Name: "Id", DataType: "int", IsPrimaryKey: true, IsNullable: false, IsIdentity: true, Description: "The identifier"},
			{Name: "Name", DataType: "varchar", IsPrimaryKey: false, IsNullable: false, IsIdentity: false, Description: "Name of the customer"},
			{Name: "Address", DataType: "varchar", IsPrimaryKey: false, IsNullable: true, IsIdentity: false, Description: "Name of the customer"},
			{Name: "ZipCode", DataType: "varchar", IsPrimaryKey: false, IsNullable: true, IsIdentity: false, Description: "Name of the customer"},
			{Name: "Budget", DataType: "decimal", IsPrimaryKey: false, IsNullable: true, IsIdentity: false, Description: "Name of the customer"},
		},
	}

	return e
}
