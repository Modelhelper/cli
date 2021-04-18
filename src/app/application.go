package app

import (
	"fmt"
	"modelhelper/cli/config"
	"os"
	"path/filepath"
	"strings"

	"runtime"

	"github.com/spf13/cobra"
)

type Initializer interface {
	Initialize() error
}

// Version shows the current application version
var Version = "3.0.0"

// Logo returns the logo to be printed on root command
func Logo() string {
	var logo = `
888b     d888               888          888 888    888          888                           
8888b   d8888               888          888 888    888          888                           
88888b.d88888               888          888 888    888          888                           
888Y88888P888  .d88b.   .d88888  .d88b.  888 8888888888  .d88b.  888 88888b.   .d88b.  888d888 
888 Y888P 888 d88""88b d88" 888 d8P  Y8b 888 888    888 d8P  Y8b 888 888 "88b d8P  Y8b 888P"   
888  Y8P  888 888  888 888  888 88888888 888 888    888 88888888 888 888  888 88888888 888     
888   "   888 Y88..88P Y88b 888 Y8b.     888 888    888 Y8b.     888 888 d88P Y8b.     888     
888       888  "Y88P"   "Y88888  "Y8888  888 888    888  "Y8888  888 88888P"   "Y8888  888     
                                                                     888                       
                                                                     888                       
                                                                     888           CLI v%v             
`
	return fmt.Sprintf(logo, Version)
}

// Info returns information about this application
func Info() string {
	infoElement := `
  Code
  ModelHelper CLI is a Command Line Interface tool to generate code based on an input source
  like a database table, REST api endpoint, a GraphQL endpoint or a proto file.
  
  Templates
  You can create your own templates based on Golang template ... each template is specified in a
  yaml- file and placed in a folder structure.

  Data
  This CLI can also help you understand database tables and perform some database tasks
  It works with MS SQL and Postgres. 

  Other input sources
  An input source can be either a database table or a set of tables. But it can also be a REST endpoint or graphql
  endpoint

  Application
  ------------
  Name:           ModelHelper 
  Version:        %v
  Location:       %v
  Environment:    %v
  Architecture:   %v
  Compiler:       %v
  Language:       go (version: %v)
  
  Drivers
  ------------
  MS SQL:         github.com/denisenkom/go-mssqldb
  RabbitMQ:       https://github.com/streadway/amqp		

  Environments
  ------------
  

  Config
  ------------
  Location:       %v
`
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	// user := "Hans-Petter Eitvet" // nødvendig? - lagres i så fall i config
	gos := runtime.GOOS
	gar := runtime.GOARCH
	gv := runtime.Version()
	gc := runtime.Compiler

	cl := "USER/.modelhelper"

	return fmt.Sprintf(infoElement, Version, exPath, gos, gar, gc, gv, cl)
}

type CommandItem struct {
	Name           string
	Short          string
	HasSubCommands bool
	Aliases        string
}

func CommandInfo(cmd *cobra.Command) string {

	cl, err := Commands(cmd)
	if err != nil {
		fmt.Printf("Error when fetching commands: %e", err)
	}

	var sb strings.Builder

	for _, ci := range cl {

		name := ci.Name

		if len(ci.Aliases) > 0 {
			name = ci.Name + " [" + ci.Aliases + "]"
		}

		fmt.Fprintf(&sb, "  %s:%s%s\n", name, " ", ci.Short)
	}

	return fmt.Sprintf(`
  Commands
  ------------

%v
`, sb.String())
}

// Commands returns a list of available commands
func Commands(cmd *cobra.Command) ([]CommandItem, error) {

	list := []CommandItem{}

	for _, c := range cmd.Commands() {

		if c.Hidden {
			continue
		}
		alias, sep := "", ""

		if len(c.Aliases) > 0 {
			for i, a := range c.Aliases {
				if i > 0 {
					sep = ", "
				}
				alias += sep + a
			}
		}

		item := CommandItem{
			c.Name(),
			c.Short,
			true,
			alias,
		}

		list = append(list, item)
	}

	return list, nil
}

func UsageInfo() string {
	return `
  Usage
  ------------

  'mh [command] [subcommand] [args] [flags]'
  `
}

func TemplateFolder(templateLocation string) string {
	var tl = ""
	if len(templateLocation) > 2 && templateLocation[0] == '.' {
		tl = filepath.Join(config.Location(), templateLocation[2:])
	}

	return tl
}

func Initialize(initializer Initializer) error {
	err := initializer.Initialize()
	return err
}

// func InitializeConfiguration() {
// 	// ConfigFolder Does not exists..
// 	rootFolder := ConfigFolder()

// 	fmt.Println("Initializing the ModelHelper configuration")

// 	err := os.Mkdir(rootFolder, os.ModeDir)
// 	if err != nil {
// 		panic(err)
// 	}

// 	f, err := os.Create(rootFolder + "/config.yaml")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()

// 	cfg := config.New()
// 	fmt.Println(cfg)
// 	_, err = f.WriteString("cfg")
// 	if err != nil {
// 		panic(err)
// 	}
// }
