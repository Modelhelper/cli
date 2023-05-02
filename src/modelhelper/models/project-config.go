package models

type ProjectConfig struct {
	Version       string            `yaml:"version"`
	Name          string            `yaml:"name"`
	Language      string            `yaml:"language"`
	Description   string            `yaml:"description"`
	OwnerName     string            `yaml:"ownerName,omitempty"`
	DefaultKey    string            `yaml:"defaultKey,omitempty"`
	DefaultSource string            `yaml:"defaultSource,omitempty"`
	RootNamespace string            `yaml:"rootNamespace,omitempty"`
	Setup         map[string]Key    `yaml:"setup,omitempty"`
	Inject        map[string]Inject `yaml:"inject,omitempty"`
	Options       map[string]string `yaml:"options,omitempty"`
	Custom        interface{}       `yaml:"custom,omitempty"`
	UseHeader     bool              `yaml:"useFileHeader"`
	Header        string            `yaml:"header,omitempty"`
	Locations     map[string]string `yaml:"locations"`
	Directory     string
	// Code          map[string]Code   `yaml:"code,omitempty"`
	// Connections   map[string]Connection `yaml:"connections,omitempty"`
}
