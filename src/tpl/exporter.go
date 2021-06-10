package tpl

import "fmt"

// This is a candidate to use standard go - check out writer
type Exporter interface {
	Export(b []byte) error
}

type ScreenExporter struct{}

func (e *ScreenExporter) Write(b []byte) (int, error) {
	fmt.Println(string(b))
	return len(b), nil
}

type DirectoryExporter struct {
	Filename  string
	Overwrite bool
}

func (e *DirectoryExporter) Write(b []byte) (int, error) {
	fmt.Println(string(b))
	return len(b), nil
}
