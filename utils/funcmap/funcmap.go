package funcmap

import (
	"sort"

	"github.com/Modelhelper/cli/v3/modelhelper/models"
	casing "github.com/Modelhelper/cli/v3/utils/text"
	"text/template"
)

func FullFuncMap(dt, ntd map[string]string) template.FuncMap {
	return FuncMap(StandardMap(), StringMap(), DatatypeMap(dt, ntd))
}

func SimpleFuncMap() template.FuncMap {
	return FuncMap(StandardMap(), StringMap())
}

func FuncMap(flist ...template.FuncMap) template.FuncMap {
	m := make(template.FuncMap)

	for _, list := range flist {
		for key, val := range list {

			m[key] = val
		}
	}

	return m
}

func StandardMap() template.FuncMap {
	return template.FuncMap{
		"increment": Incrementer,
	}

}

func Incrementer(i int) int {
	return i + 1
}

func StringMap() template.FuncMap {
	return template.FuncMap{
		"plural":   casing.PluralForm,
		"singular": casing.SingularForm,
		// "datatype":  dataTypeConverter,
		"lower":    casing.LowerCase,
		"upper":    casing.UpperCase,
		"words":    casing.AsWords,
		"sentence": casing.AsSentence,
		"snake":    casing.SnakeCase,
		"macro":    casing.MacroCase,
		"train":    casing.TrainCase,
		"kebab":    casing.KebabCase,
		"dot":      casing.DotCase,
		"title":    casing.TitleCase,
		"pascal":   casing.PascalCase,
		"camel":    casing.CamelCase,
		// "nullable":  nullableDatatype,
		// "datatypeN": dataTypeWithNullcheck,
		"append": casing.AddWord,
	}

}

func DatatypeMap(dt, ndt map[string]string) map[string]interface{} {
	m := make(map[string]interface{})

	nonull := func(input string) string {
		val, f := dt[input]

		if !f {
			return input
		}

		return val
	}

	null := func(isNullable bool, input string) string {
		if isNullable {
			val, f := ndt[input]

			if !f {
				return input
			}

			return val
		} else {
			return nonull(input)
		}
	}

	m["datatype"] = nonull
	m["datatypeN"] = null
	m["nullable"] = null

	return m
}

// LanguageMap returns the template functions that depend on the language definition of the template
//
//	imports  - the sorted, unique imports the columns need: {{ range imports .Columns }}
//	safe     - adds an underscore to a name that is reserved in the language: {{ .Name | camel | safe }}
//
// lang can be nil, imports then returns nothing and safe returns the name as it is
func LanguageMap(lang *models.LanguageDefinition) template.FuncMap {
	reserved := make(map[string]bool)
	if lang != nil {
		for _, w := range lang.ReservedWords {
			reserved[w] = true
		}
	}

	return template.FuncMap{
		"imports": func(columns []models.EntityColumnModel) []string {
			unique := make(map[string]bool)
			if lang != nil {
				for _, col := range columns {
					dt, found := lang.DataTypes[col.DataType]
					if !found {
						continue
					}

					list := dt.Imports
					if col.IsNullable {
						list = dt.NullableImports
					}

					for _, imp := range list {
						unique[imp] = true
					}
				}
			}

			out := make([]string, 0, len(unique))
			for imp := range unique {
				out = append(out, imp)
			}
			sort.Strings(out)

			return out
		},
		"safe": func(name string) string {
			if reserved[name] {
				return name + "_"
			}

			return name
		},
	}
}
