package path

import (
	"os"
)

func CurrentDirectory() string {
	cd, _ := os.Getwd()

	return cd
}
