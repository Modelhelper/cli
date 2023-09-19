package path_test

import (
	"modelhelper/cli/modelhelper/constants"
	"modelhelper/cli/utils/path"
	"path/filepath"
	"testing"
)

func Test_BasePath_Returns_correct_dir(t *testing.T) {
	cur := path.CurrentDirectory()
	pathToTest := filepath.Join(cur, "mock", "src", "some", "folder")

	expected := filepath.Join(cur, "mock", "src")
	actual, _ := path.FindBaseDirFromFoldername(pathToTest, constants.ProjectRootFolderName)

	if expected != actual {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
