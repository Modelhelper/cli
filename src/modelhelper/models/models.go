package models

import "time"

type ProjectConfig struct {
	Version       string                `yaml:"version"`
	Name          string                `yaml:"name"`
	Language      string                `yaml:"language"`
	Description   string                `yaml:"description"`
	DefaultSource string                `yaml:"defaultSource,omitempty"`
	DefaultKey    string                `yaml:"defaultKey,omitempty"`
	Connections   map[string]Connection `yaml:"connections,omitempty"`
	Code          map[string]Code       `yaml:"code,omitempty"`
	OwnerName     string                `yaml:"ownerName,omitempty"`
	Options       map[string]string     `yaml:"options,omitempty"`
	Custom        interface{}           `yaml:"custom,omitempty"`
	Header        string                `yaml:"header,omitempty"`
}

type Config struct {

	// ConfigVersion gets the version that this configuration file is using.
	ConfigVersion     string                `json:"configVersion" yaml:"configVersion"`
	AppVersion        string                `json:"appVersion" yaml:"appVersion"`
	Connections       map[string]Connection `json:"connections" yaml:"connections"`
	DefaultConnection string                `json:"defaultConnection" yaml:"defaultConnection"`
	DefaultEditor     string                `json:"editor" yaml:"editor"`
	Developer         Developer             `json:"developer" yaml:"developer"`
	Port              int                   `json:"port" yaml:"port"`
	Code              Code                  `json:"code" yaml:"code"`
	Templates         struct {
		Location string `json:"location" yaml:"location"`
	} `json:"templates" yaml:"templates"`
	Languages struct {
		Definitions string `json:"definitions" yaml:"definitions"`
	} `json:"languages" yaml:"languages"`
	Logging struct {
		Enabled bool `json:"enabled" yaml:"enabled"`
	} `json:"logging" yaml:"logging"`
}

type Developer struct {
	Name          string `json:"name" yaml:"name"`
	Email         string `json:"email" yaml:"email"`
	GitHubAccount string `json:"github" yaml:"github"`
}

type Connection struct {
	Name             string                     `json:"name" yaml:"name"`
	Description      string                     `json:"description" yaml:"description,omitempty"`
	ConnectionString string                     `json:"connectionString" yaml:"connectionString"`
	Schema           string                     `json:"schema" yaml:"schema"`
	Database         string                     `json:"database,omitempty" yaml:"database,omitempty"`
	Server           string                     `json:"server,omitempty" yaml:"server,omitempty"`
	Type             string                     `json:"type" yaml:"type"`
	Port             int                        `json:"port,omitempty" yaml:"port,omitempty"`
	Entities         []string                   `json:"entities,omitempty" yaml:"entities,omitempty"`
	Groups           map[string]ConnectionGroup `json:"groups,omitempty" yaml:"groups,omitempty"`
	Options          map[string]interface{}     `json:"options,omitempty" yaml:"options,omitempty"`
	Synonyms         map[string]string          `json:"synonyms,omitempty" yaml:"synonyms,omitempty"`
}

// should be renamed
// should this be in the input source package, since it's shared among project, config and other input sources
type ConnectionGroup struct {
	Items   []string               `json:"items" yaml:"items"`
	Options map[string]interface{} `json:"options" yaml:"options"`
}

type Synonym struct {
	Name string
}

type LanguageDefinition struct {
	Version        string              `json:"version" yaml:"version"`
	Language       string              `json:"language" yaml:"language"`
	DataTypes      map[string]Datatype `json:"datatypes" yaml:"datatypes"`
	DefaultImports []string            `json:"defaultImports" yaml:"defaultImports"`
	Keys           map[string]Key      `json:"keys" yaml:"keys"`
	Inject         map[string]Inject   `json:"inject" yaml:"inject"`
	Global         Global              `json:"global" yaml:"global"`
	Short          string              `json:"short" yaml:"short"`
	Description    string              `json:"description" yaml:"description"`
	Path           string
	// CanInject                 bool                       `json:"canInject" yaml:"canInject"`
	// UsesNamespace             bool                       `json:"usesNamespace" yaml:"usesNamespace"`
	// ModuleLevelVariablePrefix string                     `json:"moduleLevelVariablePrefix" yaml:"moduleLevelVariablePrefix"`
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

type CodeTemplate struct {
	// InjectKey       string
	// LanguageVersion string
	// Scope           TemplateScope
	// Name        string         `yaml:"name"`
	Version     string   `yaml:"version"`
	Language    string   `yaml:"language"`
	Key         string   `yaml:"key"`
	Type        string   `yaml:"type"`
	Description string   `yaml:"description"`
	Short       string   `yaml:"short"`
	Tags        []string `yaml:"tags"`
	Groups      []string `yaml:"groups"`
	FileName    string   `yaml:"fileName"`
	Model       string   `yaml:"model"`
	Body        string   `yaml:"body"`

	TemplateFilePath string
	// Export      TemplateExport `yaml:"export"`
}

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
}
type EntityModel struct {
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
}

type EntityRelationModel struct {
	IsSelfJoin        bool
	RelatedColumn     EntityColumnProps // this is either the child or parent in the relation
	OwnerColumn       EntityColumnProps // this is always the current entity
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
	// Columns            []EntityColumnImportModel
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
}

type CodeGeneratorStatistics struct {
	FilesExported    int
	TemplatesUsed    int
	EntitiesUsed     int
	SnippetsInserted int
	FilesCreated     int
	SnippetsCreated  int
	Chars            int
	Lines            int
	Words            int
	Duration         time.Duration
	TimeSaved        int
}

type CodeGeneratorResult struct {
	Statistics CodeGeneratorStatistics
	Body       []byte
	FileName   string
	Path       string
}

type CodeFileResult struct {
	Filename string
	FilePath string
	Result   *CodeGeneratorResult

	Destination  string
	Exists       bool
	Overwrite    bool
	ExistingBody []byte
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
