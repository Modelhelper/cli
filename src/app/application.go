package app

import (
	"fmt"
	"modelhelper/cli/config"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/ui"
	"os"
	"os/user"
	"path/filepath"

	"runtime"

	"github.com/gookit/color"
)

func NewModelhelperCli() modelhelper.AppService {
	return &Application{}
}

type Application struct {
	Version       string
	Configuration *modelhelper.Config
	ProjectPath   string
	IsBeta        bool
	Versions      map[string]string
}

func New() *Application {
	a := Application{}

	a.Version = version
	a.IsBeta = isBeta
	// a.Configuration = config.Load()

	a.Versions = make(map[string]string)
	a.Versions["config"] = "3.0"
	a.Versions["template"] = "3.0"
	a.Versions["project"] = "3.0"
	a.Versions["language"] = "3.0"

	return &a
}

func (a *Application) LoadConfig() *modelhelper.Config {
	loader := config.NewConfigLoader()
	cfg, _ := loader.Load()
	return cfg
}

func SetConfig(config modelhelper.Config) {
	Configuration = &config
}

var Configuration *modelhelper.Config

// version shows the current application version
var version = "3.0.0-beta2"
var isBeta = true

// Logo returns the logo to be printed on root command
func (app *Application) Logo() string {
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
                                                                     888   CLI v%v             
`
	return fmt.Sprintf(logo, version)
}

func Version() string {
	return version
}

// Info returns information about this application
func (app *Application) About() string {
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
  Name:           mh 
  Version:        %v
  Location:       %v
  Environment:    %v
  Architecture:   %v
  Compiler:       %v
  Language:       go (version: %v)
  
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

	cl := config.Location()

	return fmt.Sprintf(infoElement, version, exPath, gos, gar, gc, gv, cl)
}

func (app *Application) IsInitialized() bool {

	return config.LocationExists()
}

func UsageInfo() string {
	return `
	Usage
	------------
	
	'mh [command] [subcommand] [args] [flags]'
	`
}

func TemplateFolder(templateLocation string) string {
	if len(templateLocation) > 2 && templateLocation[0] == '.' {
		return filepath.Join(config.Location(), templateLocation[2:])
	} else {
		return templateLocation
	}

}

func (a *Application) PrintWelcomeMessage() {
	color.Green.Print(a.Logo())
	fmt.Println(Welcome())

	if isBeta {
		color.Red.Println(BetaWarning())
	}

	// fmt.Println("ModelHelper needs answers to a few questions do you wish to continue? (Y/n) ")
}

func PromptForContinue() bool {

	fmt.Printf("\nModelHelper will need to create a configuration file.")
	return ui.PromptForYesNo("Do you want to continue [Y/n]? ", "Y")
}

func Welcome() string {
	usersName := ""

	usr, err := user.Current()
	if err == nil && usr != nil && len(usr.Name) > 0 {
		usersName = fmt.Sprintf(", %s", usr.Name)
	}

	return fmt.Sprintf(`
Welcome%s to ModelHelper CLI v.%s

Code
ModelHelper is a CLI tool to generate code based on input sources
like a database table

Templates
Templates are made with the Golang template language. Each template is specified in a
yaml- file and placed in a folder structure.

Data
Understand MS SQL tables and perform some database tasks.



`, usersName, version)
	/*
	   Other input sources
	   An input source can be either a database table or a set of tables. But it can also be a REST endpoint or graphql
	   endpoint
	*/
}

func BetaWarning() string {
	return `
██╗    ██╗ █████╗ ██████╗ ███╗   ██╗██╗███╗   ██╗ ██████╗ 
██║    ██║██╔══██╗██╔══██╗████╗  ██║██║████╗  ██║██╔════╝ 
██║ █╗ ██║███████║██████╔╝██╔██╗ ██║██║██╔██╗ ██║██║  ███╗
██║███╗██║██╔══██║██╔══██╗██║╚██╗██║██║██║╚██╗██║██║   ██║
╚███╔███╔╝██║  ██║██║  ██║██║ ╚████║██║██║ ╚████║╚██████╔╝
 ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝╚═╝  ╚═══╝ ╚═════╝ 
                                                          	
This version is still in beta, expect error and strange logic.
The codebase is also chanching rapidly (not every day, but almost).

TODO:
- The command api is not finalized and may change
- Missing a few templates
- Language definitions is not done
- Optimization

Use as is and feel free to throw in any issues you may encounter,
that will help me a lot to move the code to a final state.

(any nice words is also welcome :-)
	`
}
