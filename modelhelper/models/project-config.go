package models

type ProjectConfig struct {
	Version        string                          `yaml:"version"`
	Name           string                          `yaml:"name"`
	Language       string                          `yaml:"language"`
	Description    string                          `yaml:"description"`
	OwnerName      string                          `yaml:"ownerName,omitempty"`
	DefaultKey     string                          `yaml:"defaultKey,omitempty"`
	DefaultSource  string                          `yaml:"defaultSource,omitempty"`
	RootNamespace  string                          `yaml:"rootNamespace,omitempty"`
	Features       *FeatureSet                     `yaml:"features,omitempty"`
	CustomFeatures map[string]CommonProjectFeature `yaml:"customFeatures,omitempty"`
	Setup          map[string]Key                  `yaml:"setup,omitempty"`
	Inject         map[string]Inject               `yaml:"inject,omitempty"`
	Options        map[string]string               `yaml:"options,omitempty"`
	Custom         interface{}                     `yaml:"custom,omitempty"`
	UseHeader      bool                            `yaml:"useFileHeader"`
	Header         string                          `yaml:"header,omitempty"`
	Locations      map[string]string               `yaml:"locations"`
	Directory      string
	// Code          map[string]Code   `yaml:"code,omitempty"`
	// Connections   map[string]Connection `yaml:"connections,omitempty"`
}

type CommonProjectFeature struct {
	Namespace    string            `yaml:"namespace"`
	Use          bool              `yaml:"use"`
	Options      map[string]string `yaml:"options,omitempty"`
	Imports      []string          `yaml:"imports,omitempty"`
	PropertyName *string           `yaml:"propertyName,omitempty"`
	Type         *string           `yaml:"type,omitempty"`
	Inject       *bool             `yaml:"inject,omitempty"`
	Library      *string           `yaml:"library,omitempty"`
}

type FeatureSet struct {
	Auth    *CommonProjectFeature `yaml:"auth,omitempty"`
	Logger  *CommonProjectFeature `yaml:"logger,omitempty"`
	Tracing *CommonProjectFeature `yaml:"tracing,omitempty"`
	Swagger *CommonProjectFeature `yaml:"swagger,omitempty"`
	Metrics *CommonProjectFeature `yaml:"metrics,omitempty"`
	Health  *CommonProjectFeature `yaml:"health,omitempty"`
	Api     *CommonProjectFeature `yaml:"api,omitempty"`
	Db      *CommonProjectFeature `yaml:"db,omitempty"`
}
