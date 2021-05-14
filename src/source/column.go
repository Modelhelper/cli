package source

type ColumnToTableRenderer struct {
	IncludeDescription bool
	Columns            *ColumnList
}

func (d *ColumnToTableRenderer) ToRows() [][]string {
	var rows [][]string

	for _, c := range *d.Columns {

		null := "false"
		if c.IsNullable {
			null = "true"
		}
		id := ""
		if c.IsIdentity {
			id = "yes"
		}

		pk := ""
		if c.IsPrimaryKey {
			pk = "PK"
		}

		fk := ""
		if c.IsForeignKey {
			fk = "FK"
		}

		r := []string{
			c.Name,
			c.DataType,
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

func (d *ColumnToTableRenderer) BuildHeader() []string {
	h := []string{"Name", "Type", "Nullable", "IsIdentity", "PK", "FK"}

	if d.IncludeDescription {
		h = append(h, "Description") //"Children", "Parents"
	}

	return h
}
