package exporter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ScreenExporter struct{}

func (e *ScreenExporter) Write(b []byte) (int, error) {
	fmt.Println(string(b))
	return len(b), nil
}

type FileExporter struct {
	Filename  string
	Overwrite bool
}

func (e *FileExporter) Write(b []byte) (int, error) {
	if e == nil || len(e.Filename) == 0 {
		return 0, nil
	}

	if _, err := os.Stat(e.Filename); err == nil && !e.Overwrite {
		err = errors.New("File exists")
		return 0, err
	} else {

		dir := filepath.Dir(e.Filename)
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0777)
			if err != nil {
				fmt.Println(err)
			}
		}

		err = ioutil.WriteFile(e.Filename, b, 0777)
		if err != nil {
			return 0, err
		}

		return len(b), nil
	}

}
