package source

// Entity represents an object in the relational database. Either a Table or a view
type Entity struct {
	Name               string
	ModelName          string
	ContextualName     string
	Type               string
	Schema             string
	Alias              string
	RowCount           int
	UsesIdentityColumn bool
	UsesDeletedColumn  bool
	DeletedColumnName  string
	Columns            []Column
	ParentRelations    []Relation
	ChildRelations     []Relation
	Indexes            []Index
	Description        string
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
