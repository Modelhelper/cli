package models

type TextTemplate struct {
	Template string
}
type ProjectTemplate struct {
	TemplateVersion string   `json:"templateVersion" yaml:"templateVersion"`
	Version         string   `yaml:"version"`
	Name            string   `yaml:"name"`
	Language        string   `yaml:"language"`
	Description     string   `yaml:"description,omitempty"`
	RootDirectory   string   `json:"rootDirectory" yaml:"rootDirectory"`
	Sources         []string `json:"sources,omitempty" yaml:"sources"`
	// {
	// 	Input       string `json:"input" yaml:"input"`
	// 	Destination string `json:"destination" yaml:"destination"`
	// }
	// Options          map[string]string        `yaml:"options,omitempty"`
	Tags             []string                 `json:"tags,omitempty" yaml:"tags,omitempty"`
	Commands         []ProjectTemplateCommand `json:"commands,omitempty" yaml:"commands,omitempty"`
	Wizard           []ProjectTemplateInput   `json:"wizard,omitempty" yaml:"wizard,omitempty"`
	TemplateFilePath string                   `json:"templateFilePath,omitempty" yaml:"templateFilePath,omitempty"`
	TemplateFileName string                   `json:"templateFileName,omitempty" yaml:"templateFileName,omitempty"`

	// Connections   map[string]Connection `yaml:"connections,omitempty"`
	// DefaultSource string                `yaml:"defaultSource,omitempty"`
}

type ProjectTemplateListOptions struct {
	GroupBy         string
	FilterLanguages []string
	FilterGroups    []string
	// FilterTypes     []string
	// FilterModels    []string
	// FilterKeys      []string
}

type ProjectTemplateCommand struct {
	Name    string `json:"name" yaml:"name"`
	Command string `json:"command" yaml:"command"`
}

type ProjectTemplateInput struct {
	Title       string                   `json:"title" yaml:"title"`
	Description string                   `json:"description" yaml:"description"`
	Type        string                   `json:"type" yaml:"type"`
	Returns     string                   `json:"returns" yaml:"returns"`
	Commands    []ProjectTemplateCommand `json:"commands" yaml:"commands"`
}

type ProjectTemplateModel struct {
	Name       string // admin API
	Prefix     string // stages
	RootName   string // stages
	Version    string // the application version
	SourcePath string // src/{ Name | kebab}
	RootPath   string
}
type ProjectTemplateCreateOptions struct {
	Name       string // admin API
	Template   string // admin API
	Prefix     string // stages
	RootName   string // stages
	Version    string // the application version
	SourcePath string // src/{ Name | kebab}
	RootPath   string
}

type ProjectSourceFile struct {
	DirectoryName string
	RelativePath  string
	FileName      string
	Content       []byte
}
