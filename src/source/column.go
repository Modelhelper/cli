package source

type ColumnToTableRenderer struct {
	IncludeDescription bool
	Columns            *ColumnList
}

func (d *ColumnToTableRenderer) ToRows() [][]string {
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

func (d *ColumnToTableRenderer) BuildHeader() []string {
	h := []string{"Name", "Type", "Nullable", "Identity", "PK", "FK"}

	if d.IncludeDescription {
		h = append(h, "Description") //"Children", "Parents"
	}

	return h
}
