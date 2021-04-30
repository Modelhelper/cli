package main

import (
	"fmt"
	"log"
	"modelhelper/cli/app"
	"modelhelper/cli/cmd"
	"modelhelper/cli/config"
	"modelhelper/cli/project"
	"modelhelper/cli/source"
	"path/filepath"
)

// "fmt"

func main() {

	testContext()
	//execute()

}

func execute() {
	a := app.Application{}
	a.Configuration = config.Load()

	rootExists := config.LocationExists()

	if rootExists == false {
		cfg := config.Load()
		a.Configuration = cfg

		err := a.Initialize()

		if err != nil {
			fmt.Println(err)
			panic(err)
		}

	} else {
		cmd.SetApplication(&a)
		cmd.Execute()
	}
}

func testConfig() {
	path := config.Location()
	path = filepath.Join(path, "config.yaml")
	cfg := config.LoadFromFile(path)

	println(cfg.DefaultConnection)
}
func testContext() {
	a := app.Application{}
	a.Configuration = config.Load()
	a.ProjectPath = "C:\\dev\\projects\\mh\\cli\\configuration\\project\\.modelhelper\\project.yaml"
	ctx := a.CreateContext()

	cn := ctx.DefaultConnection

	conn := ctx.Connections[cn]

	fmt.Println(conn)

}
func testLoadConnections() {
	e := project.Exists(project.DefaultLocation())
	path := ""

	if !e {
		path = "C:\\dev\\projects\\mh\\cli\\configuration\\project\\.modelhelper\\project.yaml"
	}
	p, err := project.Load(path)

	if err != nil {
		fmt.Println(err)
	}

	l := testLoader{}

	merged, _ := source.MergeConnections(&l, p)

	fmt.Println("Merged")

	for mk, vk := range *merged {
		fmt.Println(mk, vk.Description)
	}
}
func testLoadLang() {
	path := "C:\\dev\\projects\\mh\\cli\\configuration\\config\\.modelhelper\\languages\\cs.yaml"

	lang, err := config.LoadLanguageFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(lang)
}

type testLoader struct{}

func (t *testLoader) GetConnections() (*map[string]source.Connection, error) {
	var o = make(map[string]source.Connection)
	c := source.Connection{
		Description: "this is a connection comming from nowhere",
	}
	o["lab"] = c
	o["testing"] = c

	return &o, nil
}
