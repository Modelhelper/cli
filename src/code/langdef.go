package code

type DefinitionLoader interface {
	LoadDefinitions() (*map[string]LanguageDefinition, error)
}

type LanguageDefinition struct {
	Version                   string                     `json:"version" yaml:"version"`
	Language                  string                     `json:"language" yaml:"language"`
	DataTypes                 map[string]LangDefDataType `json:"datatypes" yaml:"datatypes"`
	DefaultImports            []string                   `json:"defaultImports" yaml:"defaultImports"`
	CanInject                 bool                       `json:"canInject" yaml:"canInject"`
	UsesNamespace             bool                       `json:"usesNamespace" yaml:"usesNamespace"`
	ModuleLevelVariablePrefix string                     `json:"moduleLevelVariablePrefix" yaml:"moduleLevelVariablePrefix"`
}

type LangDefDataType struct {
	Key                 string `json:"key" yaml:"key"`
	NotNull             string `json:"notNull" yaml:"notNull"`
	Nullable            string `json:"nullable" yaml:"nullable"`
	NullableAlternative string `json:"nullableAlternative" yaml:"nullableAlternative"`
}

type LangDefInject struct {
	Name         string   `json:"name" yaml:"name"`
	PropertyName string   `json:"propertyName" yaml:"propertyName"`
	Imports      []string `json:"imports" yaml:"imports"`
}

type LangDefKey struct {
	Postfix   string   `json:"postfix" yaml:"postfix"`
	Prefix    string   `json:"prefix" yaml:"prefix"`
	Imports   []string `json:"imports" yaml:"imports"`
	Inject    []string `json:"inject" yaml:"inject"`
	Namespace string   `json:"namespace" yaml:"namespace"`
}

//NewOOType return a language def for Object Oriented languages like C#, Java
func NewOOType() *LanguageDefinition {
	var l *LanguageDefinition

	l.CanInject = true
	l.UsesNamespace = true
	l.ModuleLevelVariablePrefix = "_"
	l.Version = "3.0"
	return l
}

//NewSimpleType creates a simple definition for no OO languages like C, Go
func NewSimpleType() *LanguageDefinition {
	var l *LanguageDefinition

	l.CanInject = false
	l.UsesNamespace = false
	l.ModuleLevelVariablePrefix = ""
	l.Version = "3.0"
	return l
}

func LoadDefFromPath(path string) (*LanguageDefinition, error) {
	return nil, nil
}
