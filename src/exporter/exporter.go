package exporter

type Exporter interface {
	Export(b []byte) error
}
