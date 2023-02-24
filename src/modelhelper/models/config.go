package models

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
		Code    []string `json:"code" yaml:"code"`
		Project []string `json:"project" yaml:"project"`
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
