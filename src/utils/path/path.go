package path

import (
	"log"
	"os"
	"strings"
)

func CurrentDirectory() string {
	cd, _ := os.Getwd()

	return cd
}

func FindBaseDirFromFoldername(startPath, foldername string) (string, bool) {
	if startPath == "" {
		startPath = CurrentDirectory()
	}
	folders := strings.Split(startPath, string(os.PathSeparator))

	for i := len(folders); i > 2; i-- {
		testPath := strings.Join(folders[:i], string(os.PathSeparator))
		files, err := os.ReadDir(testPath)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {

			if f.IsDir() && f.Name() == foldername {
				return testPath, true
			}
		}
	}

	return "", false
}
