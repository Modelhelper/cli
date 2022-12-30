package config

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/source"
	"modelhelper/cli/ui"

	"github.com/gookit/color"
)

// Initialize builds the configuration
func NewWithWizard() modelhelper.ConfigLoader {
	loader := &rootConfigLoader{}
	loader.config = &modelhelper.Config{}

	fmt.Println(ui.ConsoleTitle("\nMake the modelhelper CLI useful...\n"))
	fmt.Printf(`
The only way to make the modelhelper CLI useful is to

	- add at least one connection to the database
	- have access to a set of templates
	- define your favorite programming languages

You can do this now, or you can wait until you need code, or a better understanding of the database.

Anyways, if you choose to do this later you can always open the configuration file with:

	%s

`, color.Gray.Sprint("mh config --open --editor code #( or vim, sublime or what ever you use daily)"))

	do := ui.PromptForYesNo("Continue adding data to the configuration file [Y/n]? ", "Yes")

	if do {

		con := askForConnection()

		loader.config.Connections = make(map[string]modelhelper.Connection)
		loader.config.Connections[con.Name] = *con

		loader.config.Templates.Location = askForTemplateLocation()
		loader.config.Languages.Definitions = askForLanguageLocation()

		fmt.Println()
		loader.config.DefaultEditor = ui.PromptForString("Enter your default editor")

	}

	loader.path = Location()
	return loader // c.Save(Location())

}

func askForConnection() *modelhelper.Connection {

	con := modelhelper.Connection{}
	fmt.Printf("The next set of questions will build your first connection string to a MS SQL database.\n\n")
	fmt.Printf("Leave name empty if you want to skip the rest of the connection builder\n")

	con.Name = ui.PromptForString("The name of the connection (make it short and one word)")

	if con.Name == "" {
		return &con
	}

	con.Server = ui.PromptForString("The Server to connect to")
	con.Database = ui.PromptForString("The Database name to connect to")
	con.Schema = ui.PromptForString("The Default schema to be used")
	con.Type = "mssql"

	fmt.Printf("\nFor now, only username and password connections are supported.\n")
	fmt.Printf("The password and username are stored in the configuration file without encryption.\n")
	fmt.Printf("Do not share your config file\n")

	user := ui.PromptForString("Enter the Username")
	pwd := ui.PromptForPassword("Enter the Password")

	con.ConnectionString = source.BuildConnectionstring(con.Type, con.Server, con.Database, user, pwd)
	return &con
}
func askForTemplateLocation() string {

	fmt.Printf(`
To generate code the modelhelper CLI needs to know where to find a set of valid templates.

I have created a set of templates that you can copy from the modelhelper GitHub release page, 
or you can clone (or fork) directly from the template repository:

	%s

If cloned, please point to the 'code' folder in that repository...

	`, color.Gray.Sprint("git clone https://github.com/Modelhelper/templates.git"))

	loc := ui.PromptForString("Enter the path to where the templates are located")

	return loc
}
func askForLanguageLocation() string {

	fmt.Printf(`
To generate code the modelhelper CLI and the template needs to know where to find a set of valid language definitions.

I have created a set of language definitions that you can copy from the modelhelper GitHub release page, 
or you can clone (or fork) directly from the language definition repository:
	
	%s

If cloned, please point to the 'definitions' folder in that repository...

	`, color.Gray.Sprint("git clone https://github.com/Modelhelper/code-defintions.git"))

	loc := ui.PromptForString("Enter the path to language definitions")

	return loc
}
