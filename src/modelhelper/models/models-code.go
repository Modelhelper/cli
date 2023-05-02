package models

import "time"

type TemplateGeneratorStatistics struct {
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

type TemplateGeneratorResult struct {
	Statistics TemplateGeneratorStatistics
	Body       []byte
	// FileName   string
	// Path       string
}

type CodeGenerateResult struct {
	Files      []TemplateGeneratorFileResult
	Snippets   []TemplateGeneratorFileResult
	Statistics *TemplateGeneratorStatistics
	// Options    *CodeGeneratorOptions
}
type TemplateGeneratorFileResult struct {
	Filename          string
	FilePath          string
	Result            *TemplateGeneratorResult
	Destination       string
	IsSnippet         bool
	SnippetIdentifier string
	Exists            bool
	Overwrite         bool
	ExistingBody      []byte
}
type TemplateGeneratorSnippetResult struct {
	Identifer string
	Result    *TemplateGeneratorResult
}

type CodeGeneratorOptions struct {
	Name                string
	Custom              interface{}
	Templates           []string
	FeatureTemplates    []string
	TemplatePath        string
	CanUseTemplates     bool
	SourceItemGroups    []string
	SourceItems         []string
	ExportToScreen      bool
	ExportByLocationKey bool
	ExportPath          string
	ConnectionName      string
	ExportToClipboard   bool
	Overwrite           bool
	Relations           string
	CodeOnly            bool
	UseDemo             bool
	ConfigFilePath      string
	ProjectFilePath     string
	RunInteractively    bool
}
