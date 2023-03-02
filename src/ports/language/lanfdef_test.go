package language

import (
	"fmt"
	"testing"
)

func TestCouldLoadFiles(t *testing.T) {
	defs, err := loadInternalFiles()

	if err != nil {
		fmt.Println(err)
	}

	for k, v := range defs {
		fmt.Println(k, v.Language)
	}
}
