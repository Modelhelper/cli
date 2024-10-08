package models

import "time"

type CodeTemplateListOptions struct {
	DatabaseType    string
	GroupBy         string
	FilterTypes     []string
	FilterLanguages []string
	FilterModels    []string
	FilterKeys      []string
	FilterGroups    []string
	HideColumns     map[string]bool
	// SkipGroups      bool
	// SkipKey         bool

}
type SourceListOptions struct {
	ConnectionName  string
	GroupBy         string
	FilterTypes     []string
	FilterLanguages []string
	FilterModels    []string
	FilterKeys      []string
	FilterGroups    []string
	IsDemo          bool
}
type Code struct {
	RootNamespace          string            `yaml:"rootNamespace,omitempty"`
	OmitSourcePrefix       bool              `yaml:"omitSourcePrefix,omitempty"`
	Global                 Global            `yaml:"global"`
	Groups                 []string          `yaml:"groups"`
	Options                map[string]string `yaml:"options"`
	Keys                   map[string]Key    `yaml:"keys,omitempty"`
	Inject                 map[string]Inject `yaml:"inject,omitempty"`
	Locations              map[string]string `yaml:"locations"`
	FileHeader             string            `yaml:"header"`
	DisableNullableTypes   bool              `json:"diableNullableTypes" yaml:"diableNullableTypes"`
	UseNullableAlternative bool              `json:"useNullableAlternative" yaml:"useNullableAlternative"`
}
type Datatype struct {
	Key                 string      `json:"key" yaml:"key"`
	NotNull             string      `json:"notNull" yaml:"notNull"`
	Nullable            string      `json:"nullable" yaml:"nullable"`
	NullableAlternative string      `json:"nullableAlternative" yaml:"nullableAlternative"`
	DefaultValue        interface{} `json:"defaultValue" yaml:"defaultValue"`
}

type Inject struct {
	Service      string   `json:"service" yaml:"service"`
	Name         string   `json:"name" yaml:"name"`
	PropertyName string   `json:"propertyName" yaml:"propertyName"`
	Method       string   `json:"method" yaml:"method"`
	Imports      []string `json:"imports" yaml:"imports"`
}

type Key struct {
	Postfix   string   `json:"postfix" yaml:"postfix"`
	Prefix    string   `json:"prefix" yaml:"prefix"`
	Imports   []string `json:"imports" yaml:"imports"`
	Inject    []string `json:"inject" yaml:"inject"`
	Namespace string   `json:"namespace" yaml:"namespace"`
}

type Global struct {
	VariablePrefix  string `yaml:"variablePrefix"`
	VariablePostfix string `yaml:"variablePostfix"`
}

// CodeTemplate represent the full structure of a code template
type CodeTemplate struct {
	Name             string   //derived
	Version          string   `yaml:"version"` //Version denotes the version used for the template
	Language         string   `yaml:"language"`
	Identifier       string   `yaml:"identifier"`
	Key              string   `yaml:"key"`
	Type             string   `yaml:"type"`
	Description      string   `yaml:"description"`
	Short            string   `yaml:"short"`
	Tags             []string `yaml:"tags"`
	Features         []string `yaml:"features"`
	FileName         string   `yaml:"fileName"`
	Model            string   `yaml:"model"`
	Body             string   `yaml:"body"`
	BodyFile         string   `yaml:"bodyFile"`
	TemplateFilePath string
	// Export      TemplateExport `yaml:"export"`
}

// BasicModel represents the basic model for the code generation.
// It is used to generate the basic files for the project when the templates
// model is not set or set to basic
type BasicModel struct {
	RootNamespace             string
	Namespace                 string
	Postfix                   string
	Prefix                    string
	ModuleLevelVariablePrefix string
	Inject                    []InjectSection
	Imports                   []string
	Project                   ProjectSection
	Developer                 DeveloperSection
	Options                   map[string]string
	PageHeader                string
	Feature                   FeatureModel
	Name                      string
}

// FeatureModel represents the features that can be used in the project
// If a feature is configured in the project config file, it will be set to true
// and the code generation will use the features to generate the code
type FeatureModel struct {
	UseLogger  bool
	UseTracing bool
	UseMetrics bool
	UseHealth  bool
	UseApi     bool
	UseDb      bool
	UseAuth    bool
	UseSwagger bool
	Logger     FeatureOptions
	Tracing    FeatureOptions
	Metrics    FeatureOptions
	Health     FeatureOptions
	Api        FeatureOptions
	Db         DbFeatureOptions
	Auth       FeatureOptions
	Swagger    FeatureOptions
}

type FeatureOptions struct {
	PropertyName string
	Prefix       string
	Postfix      string
	Namespace    string
	Type         string
}

type DbFeatureOptions struct {
	PropertyName     string
	Prefix           string
	Postfix          string
	Namespace        string
	Type             string
	ConnectionMethod string
}

// NameModel represents the model for the name template
// If the templates model is set to name, the code generation will use this model
type NameModel struct {
	Name                      string
	RootNamespace             string
	Namespace                 string
	Postfix                   string
	Prefix                    string
	ModuleLevelVariablePrefix string
	Inject                    []InjectSection
	Imports                   []string
	Project                   ProjectSection
	Developer                 DeveloperSection
	Options                   map[string]string
	PageHeader                string
	Feature                   FeatureModel
}

// CustomModel represents the model for the custom template
type CustomModel struct {
	Custom                    interface{}
	RootNamespace             string
	Namespace                 string
	Postfix                   string
	Prefix                    string
	ModuleLevelVariablePrefix string
	Inject                    []InjectSection
	Imports                   []string
	Project                   ProjectSection
	Developer                 DeveloperSection
	Options                   map[string]string
	PageHeader                string
	Feature                   FeatureModel
}
type EntityModel struct {
	ConnectionType            string
	RootNamespace             string
	Namespace                 string
	Postfix                   string
	Prefix                    string
	ModuleLevelVariablePrefix string
	Inject                    []InjectSection
	Imports                   []string
	Project                   ProjectSection
	Developer                 DeveloperSection
	Options                   map[string]string
	Feature                   FeatureModel
	PageHeader                string
	Name                      string
	Schema                    string
	Type                      string
	Alias                     string
	Synonym                   string
	HasSynonym                bool
	ModelName                 string
	Description               string
	HasDescription            bool
	HasPrefix                 bool
	NameWithoutPrefix         string
	Columns                   []EntityColumnModel
	Parents                   []EntityRelationModel
	HasParents                bool
	Children                  []EntityRelationModel
	HasChildren               bool
	NonPrimaryColumns         []EntityColumnModel
	PrimaryKeys               []EntityColumnModel
	ForeignKeys               []EntityColumnModel
	UsedAsColumns             []EntityColumnModel
	UsesIdentityColumn        bool
	// NonIgnoredColumns  []EntityColumnViewModel
	// IgnoredColumns     []EntityColumnViewModel
}

type EntityListModel struct {
	RootNamespace             string
	Namespace                 string
	Postfix                   string
	Prefix                    string
	ModuleLevelVariablePrefix string
	Inject                    []InjectSection
	Imports                   []string
	Project                   ProjectSection
	Feature                   FeatureModel
	Developer                 DeveloperSection
	Options                   map[string]string
	PageHeader                string
	// special for entitylist
	Entities []EntityModel
}

type InjectSection struct {
	Name         string
	PropertyName string
}

type DeveloperSection struct {
	Name  string
	Email string
}

type ProjectSection struct {
	Name    string
	Owner   string
	Version string
	Exists  bool
}

type EntityRelationModel struct {
	IsSelfJoin        bool
	RelatedColumn     EntityColumnProps // this is either the child or parent in the relation
	OwnerColumn       EntityColumnProps // this is always the current entity
	OwnerShcema       string
	OwnerName         string
	Name              string
	Schema            string
	Type              string
	Alias             string
	Synonym           string
	HasSynonym        bool
	ModelName         string
	Description       string
	HasDescription    bool
	HasPrefix         bool
	NameWithoutPrefix string
	Columns           []EntityColumnModel
	NonPrimaryColumns []EntityColumnModel
	PrimaryKeys       []EntityColumnModel
	ForeignKeys       []EntityColumnModel
	UsedAsColumns     []EntityColumnModel
	// NonIgnoredColumns  []EntityColumnImportModel
	// IgnoredColumns     []EntityColumnImportModel
	// PrimaryKeys        []EntityColumnImportModel
	// ForeignKeys        []EntityColumnImportModel
	// UsedAsColumns      []EntityColumnImportModel
	UsesIdentityColumn bool
}

type EntityColumnProps struct {
	Name       string
	DataType   string
	IsNullable bool
}
type EntityColumnModel struct {
	Description       string
	IsForeignKey      bool
	IsPrimaryKey      bool
	IsIdentity        bool
	IsNullable        bool
	IsIgnored         bool
	IsDeletedMarker   bool
	IsCreatedDate     bool
	IsCreatedByUser   bool
	IsModifiedDate    bool
	IsModifiedByUser  bool
	HasPrefix         bool
	HasDescription    bool
	Name              string
	NameWithoutPrefix string
	Collation         string
	ReferencesColumn  string
	ReferencesTable   string
	DataType          string
	Length            int
	Precision         int
	Scale             int
	UseLength         bool
	UsePrecision      bool
	ForCreate         string
	IsLast            bool
	IsFirst           bool
}

type EntityList []Entity
type ColumnList []Column
type RelationList []Relation
type IndexList []Index

// Entity represents an object in the relational database. Either a Table or a view
type Entity struct {
	Name                string `json:"name" yaml:"name"`
	ModelName           string `json:"modelName" yaml:"modelName"`
	ContextualName      string `json:"contextualName" yaml:"contextualName"`
	Type                string `json:"type" yaml:"type"`
	Schema              string `json:"schema" yaml:"schema"`
	Alias               string `json:"alias" yaml:"alias"`
	Synonym             string
	HasSynonym          bool
	RowCount            int
	UsesIdentityColumn  bool
	UsesDeletedColumn   bool
	DeletedColumnName   string
	Columns             ColumnList
	ParentRelations     []Relation
	ChildRelations      []Relation
	Indexes             []Index
	Description         string
	ParentRelationCount int
	ChildRelationCount  int
	ColumnCount         int
	IdentityColumnCount int
	NullableColumnCount int
	IsVersioned         bool
	IsHistory           bool
	HistoryTable        string
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
	ColumnIndex      int
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

	ForCreate string
}

// Index represents the index of a table
type Index struct {
	ID                      string
	Name                    string
	Size                    float32
	AvgFragmentationPercent float32
	IsClustered             bool
	IsPrimaryKey            bool
	IsUnique                bool
	AvgPageSpacePercent     float32
	AvgRecordSize           float32
	Rows                    float32
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
	HasSynonym          bool
	Synonym             string
	Columns             []Column
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

type EntityRelation struct {
	IsSelfJoin bool

	ReleatedColumn   EntityColumnProps
	IncomingRelation EntityColumnProps
	OwnerColumn      EntityColumnProps
	ForeignColumn    EntityColumnProps

	// GroupIndex         int
	Name              string
	Schema            string
	Type              string
	Alias             string
	Description       string
	HasDescription    bool
	HasPrefix         bool
	NameWithoutPrefix string
	// Columns            []EntityColumnImportModel
	// NonIgnoredColumns  []EntityColumnImportModel
	// IgnoredColumns     []EntityColumnImportModel
	// PrimaryKeys        []EntityColumnImportModel
	// ForeignKeys        []EntityColumnImportModel
	// UsedAsColumns      []EntityColumnImportModel
	UsesIdentityColumn bool
}

type EntityImportModel struct {
	Code              CodeImportModel
	Options           map[string]string
	Name              string
	Schema            string
	Type              string
	RowCount          int
	Created           string
	Alias             string
	Description       string
	HasDescription    bool
	HasPrefix         bool
	NameWithoutPrefix string
	Columns           []EntityColumnImportModel
	NonPrimaryColumns []EntityColumnImportModel
	Parents           []EntityRelation
	Children          []EntityRelation
	// ModelName          string
	// ContextualName     string
	NonIgnoredColumns  []EntityColumnImportModel
	IgnoredColumns     []EntityColumnImportModel
	PrimaryKeys        []EntityColumnImportModel
	ForeignKeys        []EntityColumnImportModel
	UsedAsColumns      []EntityColumnImportModel
	UsesIdentityColumn bool
}

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
	Imports     []string
}

type CodeInjectImportModel struct {
	Key          string
	Name         string
	TemplateKeys []string
	Group        string
	PropertyName string
	Interface    string
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
	HasPrefix        bool
	HasDescription   bool
	Name             string
	// PropertyName      string
	// ContextualName    string
	NameWithoutPrefix string
	Collation         string
	ReferencesColumn  string
	ReferencesTable   string

	DataType string
	DbType   string

	Length    int
	Precision int
	Scale     int

	UseLength      bool
	UsePrecision   bool
	UseInViewModel bool
	ForCreate      string
}
