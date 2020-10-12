package main

import "modelhelper/cli/cmd"

func main() {

	rootExists := modelhelper.ConfigFolderExists()

	if rootExists == false {
		modelhelper.InitializeConfiguration()
	} else {
		cmd.Execute()
	}

}
