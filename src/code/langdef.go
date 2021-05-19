package code

type LanguageDefinition struct {
	Version        string                     `json:"version" yaml:"version"`
	Language       string                     `json:"language" yaml:"language"`
	DataTypes      map[string]LangDefDataType `json:"datatypes" yaml:"datatypes"`
	DefaultImports []string                   `json:"defaultImports" yaml:"defaultImports"`
	// CanInject                 bool                       `json:"canInject" yaml:"canInject"`
	// UsesNamespace             bool                       `json:"usesNamespace" yaml:"usesNamespace"`
	// ModuleLevelVariablePrefix string                     `json:"moduleLevelVariablePrefix" yaml:"moduleLevelVariablePrefix"`
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

func LoadFromPath(path string) (*LanguageDefinition, error) {
	return nil, nil
}
