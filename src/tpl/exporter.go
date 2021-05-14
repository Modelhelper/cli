package tpl

import "fmt"

// This is a candidate to use standard go - check out writer
type Exporter interface {
	Export(b []byte) error
}

type ScreenExporter struct{}

func (e *ScreenExporter) Export(b []byte) error {
	fmt.Println(string(b))
	return nil
}

type DirectoryExporter struct {
	Directory string
}

func (e *DirectoryExporter) Export(b []byte) error {
	fmt.Println(string(b))
	return nil
}
