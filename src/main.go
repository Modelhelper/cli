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
	"os"
	"path/filepath"
	"time"

	"github.com/olekukonko/ts"
	"gopkg.in/yaml.v3"
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
		// cmd.SetApplication(&a)
		cmd.Execute()
	}
}

func checkIdent() {
	app.PrintWelcomeMessage()
}

type Some struct {
	Version string `yaml: "version"`
	Code    Code   `yaml: "code"`
}

type Code struct {
	Keys map[string]Thing `yaml: "keys"`
}
type Thing struct {
	Postfix   string   `yaml: "postfix"`
	Namespace string   `yaml: "namespace"`
	Imports   []string `yaml: "imports,omitempty"`
}

func openProject() {
	path := "C:\\dev\\projects\\mh\\cli\\configuration\\project\\.modelhelper\\project.yaml"

	p, err := project.Load(path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ProjectLoader")
	for k, v := range p.Code.Keys {
		fmt.Println(k, v.NameSpace, v.Postfix)
	}

	dat, e := ioutil.ReadFile(path)
	if e != nil {
		log.Fatalf("cannot load file: %v", e)
	}

	var p2 *Some
	err = yaml.Unmarshal(dat, &p2)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)

	}

	fmt.Println("Unmarshal ")
	for k, v := range p2.Code.Keys {
		fmt.Println(k, v.Namespace, v.Postfix)
	}
	// x := viper.ReadConfig()

	// fmt.Println(*p2)
}
func printTerminalSizes() {
	size, _ := ts.GetSize()
	fmt.Println(size.Col())  // Get Width
	fmt.Println(size.Row())  // Get Height
	fmt.Println(size.PosX()) // Get X position
	fmt.Println(size.PosY()) // Get Y position
	//
}

func testEnvVar() {

	// (\%)(?:(?=(\\?))\2.)*?\1
	// (\%.*?\%)
	// (?<=\%)(.*?)(?=\%)
	// %SECRET:CONN_DB%
	// %ENV:CONN_DB%
	key := "PWD_LAB"

	connStr, exist := os.LookupEnv(key)
	if exist {
		fmt.Println(connStr)
	} else {
		fmt.Println("Set password " + key)
		err := os.Setenv(key, "This is a password")
		if err != nil {
			log.Fatalln(err)
		}
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
