package tpl

type Template struct {
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

type TemplateExport struct {
	FileName string `yaml:"fileName"`
	Key      string `yaml:"key"`
}
type TemplateType struct {
	Name      string `yaml:"name"`
	CanExport bool   `yaml:"canExport"`
	IsSnippet bool   `yaml:"isSnippet"`
}

var fileTemplateType = TemplateType{Name: "file", IsSnippet: true, CanExport: false}

var (
	TemplateTypes = map[string]TemplateType{
		"file":    TemplateType{Name: "file", IsSnippet: false, CanExport: true},
		"snippet": TemplateType{Name: "snippet", IsSnippet: true, CanExport: false},
		"init":    TemplateType{Name: "init", IsSnippet: false, CanExport: false},
	}
)
