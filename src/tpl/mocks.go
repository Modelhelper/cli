package tpl

func oneMoreTest() *Template {
	t := Template{}

	t.Language = "poco"
	t.Body = testTemplate()

	return &t
}

func testBlockLvl1() *Template {

	t := Template{}

	t.Name = "classname"
	t.Body = `{{ .Name }}`

	return &t
}

func testBlockLvl2() *Template {

	t := Template{}

	t.Name = "namespace"
	t.Body = `namespace {{ .NameSpace }}`

	return &t
}

func testTemplate() string {
	return `
// code her - indent by SPACE not TAB
{{- $model := index .Code.Types "model" }}
using System;
{{- range .Code.Imports }}
{{.}}
{{- end }}


//namespace from template
namespace {{ $model.NameSpace }} 
{	
	{{ template "block0" . }}
	public class {{ .Name }}{{ $model.NamePostfix}} 
	{

	}
}

{{.}}
`
}
