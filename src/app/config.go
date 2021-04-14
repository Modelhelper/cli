package app

type Config struct {
	ConfigVersion string
	AppVersion    string
	Sources       map[string]ConfigSource //[]Source
	DefaultSource string

	Logging struct {
		Enabled bool
	}
}

type ConfigSource struct {
	Name             string
	ConnectionString string
	Schema           string
	Type             string
}
