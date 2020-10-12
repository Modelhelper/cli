package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"runtime"

	"github.com/spf13/cobra"
)

// Version shows the current application version
var Version = "2.0.0"

// Logo returns the logo to be printed on root command
func Logo() string {
	var logo = ``
	return fmt.Sprintf(logo, Version)
}

// Info returns information about this application
func Info() string {
	infoElement := `

  ModelHelper CLI is a Command Line Interface tool to generate code based on datasource
  
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

	cl := "USER/.patolab"

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
