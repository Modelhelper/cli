package models

type ProjectConfig struct {
	Version     string            `yaml:"version"`
	Name        string            `yaml:"name"`
	Language    string            `yaml:"language"`
	Description string            `yaml:"description"`
	DefaultKey  string            `yaml:"defaultKey,omitempty"`
	Code        map[string]Code   `yaml:"code,omitempty"`
	OwnerName   string            `yaml:"ownerName,omitempty"`
	Options     map[string]string `yaml:"options,omitempty"`
	Custom      interface{}       `yaml:"custom,omitempty"`
	Header      string            `yaml:"header,omitempty"`
	// Connections   map[string]Connection `yaml:"connections,omitempty"`
	// DefaultSource string                `yaml:"defaultSource,omitempty"`
}
