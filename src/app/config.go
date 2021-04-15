package app

type Config struct {
	ConfigVersion string
	AppVersion    string
	Sources       map[string]ConfigSource //[]Source
	DefaultSource string

	Templates struct {
		Location string
	}
	Languages struct {
		Definitions string
	}
	Logging struct {
		Enabled bool
	}
}

type ConfigSource struct {
	Name             string
	ConnectionString string
	Schema           string
	Type             string
	Groups           map[string]ConfigSourceGroup
	Options          map[string]interface{}
}

type ConfigSourceGroup struct {
	Items   []string
	Options map[string]interface{}
}

type ConfigLanguageDef struct {
	Definitions string
}
