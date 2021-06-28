package source

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// type SortTableById []Entity
type SortTableByName []Entity
type SortTableByNameDesc []Entity
type SortTableByRows []Entity
type SortTableByRowsDesc []Entity

func (a SortTableByName) Len() int           { return len(a) }
func (a SortTableByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a SortTableByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a SortTableByNameDesc) Len() int           { return len(a) }
func (a SortTableByNameDesc) Less(i, j int) bool { return a[i].Name > a[j].Name }
func (a SortTableByNameDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a SortTableByRows) Len() int           { return len(a) }
func (a SortTableByRows) Less(i, j int) bool { return a[i].RowCount < a[j].RowCount }
func (a SortTableByRows) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a SortTableByRowsDesc) Len() int           { return len(a) }
func (a SortTableByRowsDesc) Less(i, j int) bool { return a[i].RowCount > a[j].RowCount }
func (a SortTableByRowsDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type DefaultEntitiesTableRenderer []Entity
type DescriptiveEntitiesRenderer []Entity
type SimpleEntitiesRenderer []Entity

func (d *DefaultEntitiesTableRenderer) ToRows() [][]string {
	return toRows(*d, false, true)
}

func (d *DefaultEntitiesTableRenderer) BuildHeader() []string {
	return buildHeader(false, true)
}
func (d *DescriptiveEntitiesRenderer) ToRows() [][]string {
	return toRows(*d, true, false)
}

func (d *DescriptiveEntitiesRenderer) BuildHeader() []string {
	return buildHeader(true, false)
}
func (d *SimpleEntitiesRenderer) ToRows() [][]string {
	return toRows(*d, false, false)
}

func (d *SimpleEntitiesRenderer) BuildHeader() []string {
	return buildHeader(false, false)
}

func toRows(input []Entity, withDesc, withStat bool) [][]string {
	rows := [][]string{}
	for _, e := range input {
		p := message.NewPrinter(language.English)

		r := []string{
			e.Name,
			e.Schema,
		}
		if !withDesc {
			r = append(r, e.Type)
			r = append(r, e.Alias)
			r = append(r, p.Sprintf("%d", e.RowCount))
		}
		if withStat {
			r = append(r, p.Sprintf("%d", e.ColumnCount))
			r = append(r, p.Sprintf("%d", e.ParentRelationCount))
			r = append(r, p.Sprintf("%d", e.ChildRelationCount))
		}

		if withDesc {
			r = append(r, e.Description)
		}

		rows = append(rows, r)
	}

	return rows

}

func buildHeader(withDesc, withStat bool) []string {
	h := []string{"Name", "Schema"}

	if !withDesc {
		h = append(h, "Type")
		h = append(h, "Alias")
		h = append(h, "Rows")
	}

	if withStat {
		h = append(h, "Col Cnt")
		h = append(h, "P Relations")
		h = append(h, "C Relations")

	}

	if withDesc {
		h = append(h, "Description")
	}

	return h
}
