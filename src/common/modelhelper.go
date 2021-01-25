package common

import (
	"fmt"
	"modelhelper/cli/defaults"
	"os"
)

//ConfigFolder returns the root path of ModelHelper
func ConfigFolder() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s/.modelhelper", homeDir)

}

// ConfigFolderExists checks if the config folder exists
func ConfigFolderExists() bool {
	homeDir := ConfigFolder()

	if _, err := os.Stat(homeDir); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func InitializeConfiguration() {
	// ConfigFolder Does not exists..
	rootFolder := ConfigFolder()

	fmt.Println("Initializing the ModelHelper configuration")

	err := os.Mkdir(rootFolder, os.ModeDir)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(rootFolder + "/config.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(defaults.Configuration())
	if err != nil {
		panic(err)
	}
}