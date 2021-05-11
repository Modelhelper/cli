package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"modelhelper/cli/app"
	"modelhelper/cli/cmd"
	"modelhelper/cli/config"
	"modelhelper/cli/project"
	"modelhelper/cli/source"
	"modelhelper/cli/tpl"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// "fmt"

func main() {

	// testTemplateLoader()
	//testContext()
	execute()

}

func execute() {
	a := app.Application{}

	rootExists := config.LocationExists()

	if rootExists == false {

		a.Configuration = config.Load()

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
func testTemplateLoader() {

	st := time.Now()
	fileName := "C:\\Users\\hans-petter\\.modelhelper\\templates\\tutorial\\strings.yaml"
	fmt.Println("load")
	dat, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Fatalf("cannot load file: %v", e)
	}

	fmt.Println("end")
	dur := time.Since(st)

	fmt.Printf("\nDur: %v ms", dur.Milliseconds())
	it := 10000
	loadFullTemplate(dat, it)
	loadTypeTemplate(dat, it)
	// startFull := time.Now()
	// for i := 1; i < it; i++ {
	// 	loadFullTemplate(dat)
	// }
	// durationFull := time.Since(startFull)
	// fmt.Printf("\nFull template: %v iterations, time: %v ms", it, durationFull.Milliseconds())

}

func loadFullTemplate(dat []byte, iterations int) {
	strt := time.Now()
	for i := 1; i < iterations; i++ {
		var p tpl.Template

		err := yaml.Unmarshal(dat, &p)
		if err != nil {
			log.Fatalf("cannot unmarshal data: %v", err)

		}
	}
	dur := time.Since(strt)
	fmt.Printf("\nFull template: %v iterations, time: %v ms", iterations, dur.Milliseconds())

}
func loadTypeTemplate(dat []byte, iterations int) {

	strt := time.Now()
	for i := 1; i < iterations; i++ {
		var p TypeTpl

		err := yaml.Unmarshal(dat, &p)
		if err != nil {
			log.Fatalf("cannot unmarshal data: %v", err)

		}
	}
	dur := time.Since(strt)
	fmt.Printf("\nType template: %v iterations, time: %v ms", iterations, dur.Milliseconds())

}

type TypeTpl struct {
	Type string `yaml:"type"`
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
