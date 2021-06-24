package main

import (
	"fmt"
	"modelhelper/cli/app"
	"modelhelper/cli/cmd"
)

// "fmt"

func main() {

	execute()
}

func execute() {
	a := app.Application{}

	isInitialized := a.IsInitialized()

	if isInitialized == false {

		err := a.Initialize()

		if err != nil {
			fmt.Println(err)
			panic(err)
		}

	} else {
		cmd.Execute()
	}
}

// func printTerminalSizes() {
// 	size, _ := ts.GetSize()
// 	fmt.Println(size.Col())  // Get Width
// 	fmt.Println(size.Row())  // Get Height
// 	fmt.Println(size.PosX()) // Get X position
// 	fmt.Println(size.PosY()) // Get Y position
// 	//
// }
