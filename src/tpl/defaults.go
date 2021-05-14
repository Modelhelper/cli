package tpl

func DefaultModels() map[string]Description {
	models := make(map[string]Description)

	models["basic"] = Description{
		Short: "'basic' is default. This model provides data about current project, developer etc",
		Long: `Use the basic model if the template you are creating don't has the
	need for a single entity or a list of entities`,
	}

	models["entity"] = Description{
		Short: "Includes same properties as 'basic', but will also bring in an entity",
		Long: `An entity model contains information about one entity. 
Use the entity type if you have create code that need access to columns, primary keys, foreign keys
parent and child- relations.

This is the most used model type.`,
	}
	models["entities"] = Description{
		Short: "Includes same properties as 'basic', but will also bring in a list of entities",
		Long: `The entities model contains information about all entities provided by the following:
	1) --entity name1 --entity  name2. or
	2) --entity-group groupname (a group of entities must be done on a 
	   connection property in the config or a project)

With a list of entities you will be able to loop through each entity and create code thats fits in one single page

Each entity in the list access to columns, primary keys, foreign keys, parent and child- relations.`,
	}

	return models
}

func DefaultTypes() map[string]Description {
	models := make(map[string]Description)

	models["file"] = Description{
		Short: "Includes same properties as 'basic', but will also bring in an entity",
		Long: `A file template represents complete 'artifact', for C# nd Java this could be a
single class definition or anything else that fits on a page.

A file type can be exported as file
		`,
	}
	models["block"] = Description{
		Short: "A block template can be used as a reference in another block or file",
		Long: `Using blocks as building stones will make the templates more readable
		
Blocks enable better reuse of code, smaller templates and should make templates
faster to create.

Reference the block in a snippet, file or another block with the following syntax:

		{{ template "block-name" . }}

`,
	}

	return models
}
func DefaultKeys() map[string]Description {
	models := make(map[string]Description)

	gen := `The 'key' is used to connect a generic template to parts of the a specific project
or a standard setup from a config or language definition.

a set of keys for C# can be totally different from keys in go, java or other languages.
`
	models["model"] = Description{
		Short: "model",
		Long:  gen + `The 'model' key is normally used on poco and objects representing an entity`,
	}
	models["interface"] = Description{
		Short: "interface",
		Long:  gen + `interface`,
	}
	models["repo"] = Description{
		Short: "repo",
		Long:  gen + `repo`,
	}

	return models
}
