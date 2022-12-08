package model

import "context"

type ModelConverter interface {
	ToModel(ctx context.Context) interface{}
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
	Columns                   []EntityColumnViewModel
	Parents                   []EntityRelationViewModel
	Children                  []EntityRelationViewModel
	PrimaryKeys               []EntityColumnViewModel
	ForeignKeys               []EntityColumnViewModel
	UsedAsColumns             []EntityColumnViewModel
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
	Name  string
	Owner string
}

type EntityRelationViewModel struct {
	IsSelfJoin        bool
	ReleatedColumn    EntityColumnProps // this is either the child or parent in the relation
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
type EntityColumnViewModel struct {
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
