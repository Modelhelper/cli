package source

import "modelhelper/cli/modelhelper"

type SortColumnById []modelhelper.Column
type SortColumnByName []modelhelper.Column
type ColumnToTableRenderer struct {
	IncludeDescription bool
	Columns            *modelhelper.ColumnList
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
