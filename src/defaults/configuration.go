package defaults

import "fmt"

var ApplicationVersion string = "3.0"
var TemplateVersion string = "3.0"
var ProjectVersion string = "3.0"
var ConfigVersion string = "3.0"

type TemplateOptions struct {
	Language    string
	Import      string
	LocationKey string
}

func Configuration() string {
	return `
version: 3.0
appVersion: 3.0

templates:
  shared:
  global:
    location: ./templates

languages:
  definitions: ./languages

logging:
  enabled: true
`
}

func DefaultTemplateContent(options *TemplateOptions) []byte {
	tpl := `
#mandatory, this is the template version
version: %s

# mandatory. connects to the correct language definition
language: %s

# mandatory, default: file (options: file | snippet)
type: file

# what groups this template is included in
# is used to fetch templates within a group when using -tg in generate command
# optional, default is null
groups: []

# optional, default is null
tags: []

# indicates the model classes to be injected
# optional, default is none
import: %s

# can use template syntax to use information from the models
# if snippet is used, this will be the filename to inject the snippets to
# optional if ExportType = none
exportFileName: "{{ .Table.Name | UpperCamel | Singular}}.%s"

# optional, a key cannot contain .
locationKey: %s 

# the body of the template
# mandatory
body: |
  // code her - indent by SPACE not TAB
  using System;
  {{ range .Code.Imports }}
  {{ . }}
  {{ end }}
  namespace {{ .Code.Types[key] }}
  {
    // simple if
    {{ if .Code.Creator != nil}}
      // developer name: {{ .Developer }} for {{ .Company }}
    {{ end}}

    public class {{ .Name | Pascal | Singular}}{{ .Code.Types['model'].NamePostfix}}
    {
    {{ range .NonIgnoredColumns }}
      public {{ .Name | Pascal }} { get; set; }
    {{ end }}    
    }
  }
`

	return []byte(fmt.Sprintf(tpl, TemplateVersion, options.Language, options.Import, options.Language, options.LocationKey))
}

type ProjectOptions struct {
	Version     string
	Language    string
	Import      string
	LocationKey string
}

func DefaultProjectSetup(options *ProjectOptions) []byte {
	tpl := `
  #mandatory, this is the project version
  version: %s
  
  # mandatory. connects to the correct language definition
  language: %s
  
  # mandatory, default: file (options: file | snippet)
  type: file
  
  # what groups this template is included in
  # is used to fetch templates within a group when using -tg in generate command
  # optional, default is null
  groups: []
  
  # optional, default is null
  tags: []
  
  # indicates the model classes to be injected
  # optional, default is none
  import: %s
  
  # can use template syntax to use information from the models
  # if snippet is used, this will be the filename to inject the snippets to
  # optional if ExportType = none
  exportFileName: "{{ .Table.Name | UpperCamel | Singular}}.cs"
  
  # optional, a key cannot contain .
  locationKey: %s 
  
  # the body of the template
  # mandatory
  body: |
   // code her - indent by SPACE not TAB
  `

	return []byte(fmt.Sprintf(tpl, ProjectVersion, &options.Language, &options.Import, &options.LocationKey))
}
