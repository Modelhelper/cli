package source

type Source struct {
	Name             string
	ConnectionString string
	Schema           string
	Type             string
	Groups           map[string]Group
	Options          map[string]interface{}
}

// should this be in the input source package, since it's shared among project, config and other input sources
type Group struct {
	Items   []string
	Options map[string]interface{}
}

func New() *Source {
	return nil
}
