package tpl

import "fmt"

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
