package code

import (
	"fmt"
	"testing"
)

func TestCouldLoadFiles(t *testing.T) {
	defs, err := Load()

	if err != nil {
		fmt.Println(err)
	}

	for k, v := range defs {
		fmt.Println(k, v.Language)
	}
}
