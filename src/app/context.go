package app

import (
	"log"
	"modelhelper/cli/code"
	"modelhelper/cli/config"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/project"
)

// Templates         *[]tpl.Template
// Blocks            *[]tpl.Template
type Context struct {
	ProjectExists     bool
	CurrentProject    *modelhelper.ProjectConfig
	Connections       map[string]modelhelper.Connection
	DefaultConnection string
	Languages         *map[string]code.LanguageDefinition
	Options           *map[string]interface{}
	CurrentLanguage   code.LanguageDefinition
	CurrentConnection modelhelper.Connection
	IsDemo            bool
	InputConnection   string
}

func (a *Application) CreateContext() *Context {
	c := Context{}

	con := make(map[string]modelhelper.Connection)

	if a.Configuration == nil {
		a.Configuration = config.Load()
	}

	if len(a.Configuration.DefaultConnection) > 0 {
		c.DefaultConnection = a.Configuration.DefaultConnection
	}

	if a.Configuration.Connections != nil {

		for ck, cv := range a.Configuration.Connections {
			con[ck] = cv
		}
	}

	if len(a.ProjectPath) == 0 {
		a.ProjectPath = project.DefaultLocation()
	}

	c.ProjectExists = project.Exists(a.ProjectPath)

	if c.ProjectExists {
		pr, err := project.Load(a.ProjectPath)
		if err != nil {
			log.Fatalln(err)
		}
		c.CurrentProject = &pr.Config

		for pk, pv := range c.CurrentProject.Connections {
			con[pk] = pv
		}

		if len(c.CurrentProject.DefaultSource) > 0 {
			c.DefaultConnection = c.CurrentProject.DefaultSource
		}
	}

	c.Connections = con

	return &c
}
