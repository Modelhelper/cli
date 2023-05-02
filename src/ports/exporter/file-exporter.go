package exporter

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

func NewSnippetExporter(fn, id string) io.Writer {
	return &snippetExporter{fn, id}
}

type snippetExporter struct {
	fileName   string
	identifier string
}

func (e *snippetExporter) Write(b []byte) (int, error) {
	writeSnippetFile(e.identifier, string(b), e.fileName)
	return 1, nil
}

func writeSnippetFile(snippetStringIdentifier, snippetCode, filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	modifiedContent := writeSnippet(snippetStringIdentifier, snippetCode, content)

	// Write the modified content back to the file
	err = ioutil.WriteFile(filePath, modifiedContent, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}

func writeSnippet(snippetStringIdentifier, snippetCode string, content []byte) []byte {

	snippetStringIdentifier = strings.ReplaceAll(snippetStringIdentifier, "%", "")
	// Define the regular expression pattern to match the search text
	searchText := fmt.Sprintf(`%%%%%s%%%%`, snippetStringIdentifier)
	pattern := regexp.MustCompile(searchText)

	// Find the position of the search text using the regular expression
	matches := pattern.FindAllStringSubmatchIndex(string(content), -1)
	// if len(matches) < 2 {
	// 	fmt.Printf("Error: could not find \"%s\" in file\n", searchText)
	// 	return content
	// }

	for x := len(matches) - 1; x >= 0; x-- {
		value := matches[x]

		insertIndex := value[1] + 1
		cnt := []byte(string(content)[:insertIndex] + snippetCode + "\n" + string(content)[insertIndex:])
		content = cnt
		//fmt.Printf("from Pattern %d", pattern.FindStringIndex(string(content))[1]) // add 1 to move past the newline character
	}

	return content
}
