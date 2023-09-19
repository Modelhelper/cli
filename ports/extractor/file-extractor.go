package extractor

import (
	"fmt"
	"io/ioutil"
	"modelhelper/cli/modelhelper"
	"regexp"
)

type dockerStactVariableExtractor struct {
	filename string
}

// Extract implements modelhelper.TextExtractor
func (e *dockerStactVariableExtractor) Extract() []string {
	expression := `\$\{([A-Z\_].*)\}`

	content, err := ioutil.ReadFile(e.filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	pattern := regexp.MustCompile(expression)
	matches := pattern.FindAllStringSubmatch(string(content), -1)

	list := []string{}

	for _, match := range matches {
		if len(match) > 1 {
			m := match[1]
			list = append(list, m)
		}
	}

	return list
}

func NewDockerStackVariableExtractor(filename string) modelhelper.TextExtractor {
	return &dockerStactVariableExtractor{filename: filename}
}
