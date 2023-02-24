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
	FileName   string
	Path       string
}

type TemplateGeneratorFileResult struct {
	Filename string
	FilePath string
	Result   *TemplateGeneratorResult

	Destination  string
	Exists       bool
	Overwrite    bool
	ExistingBody []byte
}

type CodeGeneratorOptions struct {
	Templates         []string
	TemplateGroups    []string
	TemplatePath      string
	CanUseTemplates   bool
	EntityGroups      []string
	Entities          []string
	ExportToScreen    bool
	ExportByKey       bool
	ExportPath        string
	Connection        string
	ExportToClipboard bool
	Overwrite         bool
	Relations         string
	CodeOnly          bool
	UseDemo           bool
	ConfigFilePath    string
	ProjectFilePath   string
	RunInteractively  bool
}
