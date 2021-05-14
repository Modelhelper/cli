package app

import (
	"log"
	"modelhelper/cli/code"
	"modelhelper/cli/config"
	"modelhelper/cli/project"
	"modelhelper/cli/source"
)

type Context struct {
	ProjectExists  bool
	CurrentProject *project.Project
	// Templates         *[]tpl.Template
	// Blocks            *[]tpl.Template
	Connections       map[string]source.Connection
	DefaultConnection string
	Languages         *map[string]code.LanguageDefinition
	Options           *map[string]interface{}
	CurrentLanguage   code.LanguageDefinition
	CurrentConnection source.Connection
}

func (a *Application) CreateContext() *Context {
	c := Context{}

	con := make(map[string]source.Connection)

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
		c.CurrentProject = pr

		for pk, pv := range pr.Connections {
			con[pk] = pv
		}

		if len(pr.DefaultSource) > 0 {
			c.DefaultConnection = pr.DefaultSource
		}
	}

	c.Connections = con

	return &c
}
