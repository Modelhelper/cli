package main

import (
	"fmt"
	"os"
)

var Version string = "3.0.0"

func main() {

	var name string
	fmt.Println("What is your name?")
	fmt.Scanf("%s\n", &name)

	var age int
	fmt.Println("What is your age?")
	fmt.Scanf("%d\n", &age)

	fmt.Printf("Hello %s, your age is %d\n", name, age)

	f, err := os.Open("test.file")
	if err != nil {
		panic(err)
	}
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	fmt.Printf("File name: %v\n", fi.Name())
	fmt.Printf("Is Directory: %t\n", fi.IsDir())
	fmt.Printf("Size: %d\n", fi.Size())
	fmt.Printf("Mode: %v\n", fi.Mode())

	// check if .modelhelper exists
	// rootExists := modelhelper.ConfigFolderExists()

	// if rootExists == false {
	// 	modelhelper.InitializeConfiguration()
	// } else {
	// 	cmd.Execute()
	// }

}
