package types

type Project struct {
	Version       string
	Name          string
	DefaultSource string
	Sources       []ProjectSource
	Code          []ProjectCode
	Options       map[string]string
	Custom        interface{}
	RootPath      string
}

type ProjectSource struct {
	Name       string
	Connection string
	Schema     string
	Type       SourceType
	Mapping    []SourceMapping
	Groups     []SourceGroup
	Options    map[string]string
}

type SourceType string

const (
	mssql SourceType = "mssql"
)

type SourceMapping struct {
	Name            string
	IsIgnored       bool
	TrueValue       interface{}
	FalseValue      interface{}
	IsDeletedMarker bool
	IsCreationDate  bool
	IsModifiedDate  bool
}

type SourceGroup struct {
	Name  string
	Items []string
}
type ProjectCode struct {
	OmitSourcePrefix bool
	Global           GlobalCode
	Groups           []string
	Options          map[string]string
}

type CodeInject struct {
	Name         string
	Language     string
	PropertyName string
	Interface    string
	Namespace    string
	Method       string
	Imports      []string
}

type GlobalCode struct {
	VariablePrefix  string
	VariablePostfix string
}
type CodeKey struct {
	Name        string
	Group       string
	Path        string
	NameSpace   string
	NamePostfix string
	NamePrefix  string
}
