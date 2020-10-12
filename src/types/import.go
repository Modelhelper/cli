package types

import "time"

type CreatorImportModel struct {
	CompanyName   string
	DeveloperName string
}

type CodeImportModel struct {
	Creator               CreatorImportModel
	OmitSourcePrefix      bool
	CurrentDate           time.Time
	GlobalVariablePrefix  string
	GlobalVariablePostfix string
	CanInject             bool
	Inject                map[string]CodeInjectImportModel
	Types                 map[string]CodeTypeImportModel
	Imports               []string
	Language              string
	// Locations             []CodeLocationImportModel
}

type CodeTypeImportModel struct {
	Key         string
	NamePostfix string
	NamePrefix  string
	NameSpace   string
}

type CodeInjectImportModel struct {
	Key          string
	Name         string
	TemplateKeys []string
	Group        string
	PropertyName string
	Interface    string
}

type EntityImportModel struct {
	Code               CodeImportModel
	Options            map[string]string
	Name               string
	Schema             string
	Type               string
	RowCount           int
	Created            string
	Alias              string
	Description        string
	Columns            []EntityColumnImportModel
	Parents            []EntityRelation
	Children           []EntityRelation
	ModelName          string
	ContextualName     string
	NonIgnoredColumns  []EntityColumnImportModel
	IgnoredColumns     []EntityColumnImportModel
	PrimaryKeys        []EntityColumnImportModel
	ForeignKeys        []EntityColumnImportModel
	UsedAsColumns      []EntityColumnImportModel
	UsesIdentityColumn bool
}

type EntityColumnImportModel struct {
	Description      string
	IsForeignKey     bool
	IsPrimaryKey     bool
	IsIdentity       bool
	IsNullable       bool
	IsIgnored        bool
	IsDeletedMarker  bool
	IsCreatedDate    bool
	IsCreatedByUser  bool
	IsModifiedDate   bool
	IsModifiedByUser bool

	Name             string
	PropertyName     string
	ContextualName   string
	Collation        string
	ReferencesColumn string
	ReferencesTable  string

	DataType string
	DbType   string

	Length    int
	Precision int
	Scale     int

	UseLength      bool
	UsePrecision   bool
	UseInViewModel bool
}
