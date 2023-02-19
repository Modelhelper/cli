package modelhelper

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (d *EntityList) Rows() [][]string {
	var rows [][]string

	for _, e := range *d {
		p := message.NewPrinter(language.English)

		r := []string{
			e.Name,
			e.Schema,
			e.Alias,
			p.Sprintf("%d", e.RowCount),
			p.Sprintf("%d", e.ColumnCount),
			p.Sprintf("%d", e.ParentRelationCount),
			p.Sprintf("%d", e.ChildRelationCount),
			// strconv.Itoa(len(e.ChildRelations)),
			// strconv.Itoa(len(e.ParentRelations)),
		}

		// if withDesc {
		// 	r = append(r, e.Description)
		// }

		rows = append(rows, r)
	}

	return rows

}

func (d *EntityList) Header() []string {
	h := []string{"Name", "Schema", "Alias", "Rows", "Col Cnt", "P Relations", "C Relations"}

	// if withDesc {
	// 	h = append(h, "Description"), "Children", "Parents"
	// }

	return h
}
