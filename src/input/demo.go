package input

type DemoInput struct{}

func (server *DemoInput) Entity(name string) (*Entity, error) {

	entities := getEntities()
	//var entity Entity

	for _, e := range entities {
		if e.Name == name {
			return &e, nil
		}
	}

	return nil, nil
}
func (server *DemoInput) Entities(pattern string) (*[]Entity, error) {
	e := getEntities()
	return &e, nil
}

func getEntities() []Entity {
	e := []Entity{

		getOrderHeadTable(),
		getCustomerTable(),
	}
	return e
}

func getOrderHeadTable() Entity {
	return Entity{
		Name: "order", Schema: "dbo", Description: "This is the order table", RowCount: 1000, Alias: "o", UsesIdentityColumn: true,
		Columns: []Column{
			{Name: "Id", DataType: "int", IsPrimaryKey: true, IsNullable: false, IsIdentity: true, Description: "The identifier"},
			{Name: "Name", DataType: "varchar", IsPrimaryKey: false, IsNullable: false, IsIdentity: false, Description: "Name of the order"},
		},
	}
}
func getCustomerTable() Entity {
	var e = Entity{
		Name: "customer", Schema: "dbo",
		Description: "This is the customer table", RowCount: 1000, Alias: "c", UsesIdentityColumn: true,
		Columns: []Column{
			{Name: "Id", DataType: "int", IsPrimaryKey: true, IsNullable: false, IsIdentity: true, Description: "The identifier"},
			{Name: "Name", DataType: "varchar", IsPrimaryKey: false, IsNullable: false, IsIdentity: false, Description: "Name of the customer"},
		},
	}

	return e
}
