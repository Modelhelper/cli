package funcmap

import (
	casing "modelhelper/cli/utils/text"
	"text/template"
)

func FullFuncMap(dt, ntd map[string]string) template.FuncMap {
	return FuncMap(StringMap(), DatatypeMap(dt, ntd))
}

func SimpleFuncMap() template.FuncMap {
	return FuncMap(StringMap())
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
