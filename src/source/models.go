package source

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type EntityList []Entity
type ColumnList []Column
type RelationList []Relation
type IndexList []Index

// Entity represents an object in the relational database. Either a Table or a view
type Entity struct {
	Name               string `json:"name" yaml:"name"`
	ModelName          string `json:"modelName" yaml:"modelName"`
	ContextualName     string `json:"contextualName" yaml:"contextualName"`
	Type               string `json:"type" yaml:"type"`
	Schema             string `json:"schema" yaml:"schema"`
	Alias              string `json:"alias" yaml:"alias"`
	RowCount           int
	UsesIdentityColumn bool
	UsesDeletedColumn  bool
	DeletedColumnName  string
	Columns            ColumnList
	ParentRelations    RelationList
	ChildRelations     RelationList
	Indexes            IndexList
	Description        string
}

func (d *EntityList) ToRows() [][]string {
	var rows [][]string

	for _, e := range *d {
		p := message.NewPrinter(language.English)

		r := []string{
			e.Name,
			e.Schema,
			e.Alias,
			p.Sprintf("%d", e.RowCount),
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

func (d *EntityList) BuildHeader() []string {
	h := []string{"Name", "Schema", "Alias", "Rows"}

	// if withDesc {
	// 	h = append(h, "Description"), "Children", "Parents"
	// }

	return h
}

type EntityStat struct {
	Schema        string
	Name          string
	Description   string
	PkCount       int
	FkCount       int
	RowCount      int
	ColumnCount   int
	ChildrenCount int
	ParentCount   int
	IndexCount    int
	Size          int
}

// Column represents the column of an entity, either a table or a view
type Column struct {
	ID               int
	Name             string
	PropertyName     string
	DbType           string
	DataType         string
	Collation        string
	IsPrimaryKey     bool
	IsForeignKey     bool
	IsNullable       bool
	IsIdentity       bool
	IsIgnored        bool
	IsCreatedByUser  bool
	IsCreatedDate    bool
	IsModifiedByUser bool
	IsModifiedDate   bool
	IsDeletedMarker  bool
	Precision        int
	Scale            int
	Length           int
	UsePrecision     bool
	UseLength        bool
	UseInViewModel   bool
	IsReserved       bool
	ReferencesTable  string
	ReferencesColumn string

	Description    string
	ContextualName string
}

// Index represents the index of a table
type Index struct {
	ID                      string
	Name                    string
	Size                    int
	AvgFragmentationPercent int
	IsClustered             bool
	IsPrimaryKey            bool
	IsUnique                bool
	AvgPageSpacePercent     int
	AvgRecordSize           int
	Rows                    int
	Columns                 []IndexColumn
}

type IndexColumn struct {
	Name              string
	IsDescending      bool
	IsNullable        bool
	IsIdentity        bool
	PartitionOriginal int
}

type TableRelation struct {
	GroupIndex           int
	ConstraintName       string
	ParentColumnName     string
	ChildColumnName      string
	ParentColumnType     string
	ChildColumnType      string
	ParentColumnNullable bool
	ChildColumnNullable  bool
	IsSelfJoin           bool
}

type Relation struct {
	GroupIndex          int
	Name                string
	Schema              string
	Type                string
	SortIndex           int
	Depth               int
	Family              string
	OwnerColumnName     string
	OwnerColumnType     string
	OwnerColumnNullable bool
	ColumnName          string
	ColumnType          string
	ColumnNullable      bool
	ContraintName       string
	IsSelfJoin          bool
	// Level               int
	// FullPath            string
	// ReferenceName       string
	// ParentTableName       string
	// ReferencedTableName   string
	// ForeignColumnName     string
	// ForeignColumnType     string
	// ForeignColumnNullable bool
}

type DatabaseInformation struct {
	Version    string
	ServerName string
}
