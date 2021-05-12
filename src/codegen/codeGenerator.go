package codegen

import (
	"io/ioutil"
	"log"
	"modelhelper/cli/app"
	"modelhelper/cli/code"
	"modelhelper/cli/tpl"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/gertd/go-pluralize"
)

var _ctx *app.Context

type Result struct {
	Milliseconds int
}

type TemplateExport struct {
	FileName string `yaml:"fileName"`
	Key      string `yaml:"key"`
}
type TemplateType struct {
	Name      string `yaml:"name"`
	CanExport bool   `yaml:"canExport"`
	IsSnippet bool   `yaml:"isSnippet"`
}

var fileTemplateType = TemplateType{Name: "file", IsSnippet: true, CanExport: false}

var (
	TemplateTypes = map[string]TemplateType{
		"file":    TemplateType{Name: "file", IsSnippet: false, CanExport: true},
		"snippet": TemplateType{Name: "snippet", IsSnippet: true, CanExport: false},
		"init":    TemplateType{Name: "init", IsSnippet: false, CanExport: false},
	}
)

func Generate(model interface{}, ctx *app.Context) (*Result, error) {
	_ctx = ctx

	return nil, nil
}

func useTempdir(name string, blocks []*tpl.Template) *template.Template {

	dir := createTempDir()
	defer os.RemoveAll(dir)

	err := writeTempFiles(dir, blocks)

	if err != nil {
		return nil
	}

	pattern := filepath.Join(dir, "*")

	fm := template.FuncMap{
		"plural":    pluralForm,
		"singular":  SingularForm,
		"datatype":  dataTypeConverter,
		"lower":     lowerCase,
		"upper":     upperCase,
		"words":     asWords,
		"sentence":  asSentence,
		"snake":     snakeCase,
		"kebab":     kebabCase,
		"pascal":    pascalCase,
		"camel":     camelCase,
		"nullable":  nullableDatatype,
		"datatypeN": dataTypeWithNullcheck,
		"append":    addWord,
	}

	drivers := template.Must(template.New(name).Funcs(fm).ParseGlob(pattern))

	return drivers
}

func createTempDir() string {

	dir, err := ioutil.TempDir("", "template") //os.MkdirTemp("", "template")
	if err != nil {
		log.Fatal(err)
	}

	return dir

}

func writeTempFiles(dir string, files []*tpl.Template) error {
	for _, file := range files {
		fp := filepath.Join(dir, file.Name)
		err := ioutil.WriteFile(fp, []byte(file.Body), 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

func pluralForm(input string) string {
	pluralize := pluralize.NewClient()
	output := pluralize.Plural(input)

	return output
}
func SingularForm(input string) string {
	pluralize := pluralize.NewClient()
	output := pluralize.Singular(input)

	return output
}

func dataTypeWithNullcheck(isNullable bool, input string) string {
	return ""

	// dt := dataTypeConverter(input)

	// // if reflect.TypeOf(isNullable) == reflect.Typeof(bool)
	// if isNullable {
	// 	return nullableDatatype(dt)
	// }
	// return dt
}
func alternativeNullableDatatype(input string) string {
	dict := make(map[string]string)
	dict["int"] = "Nullable<int>"
	dict["long"] = "Nullable<long>"
	dict["bool"] = "Nullable<bool>"

	output := dict[input]

	if len(output) > 0 {
		return output
	}

	return input
}
func nullableDatatype(input string) string {
	dict := make(map[string]string)
	dict["int"] = "int?"
	dict["long"] = "long?"
	dict["bool"] = "bool?"

	output := dict[input]

	if len(output) > 0 {
		return output
	}

	return input
}

func dataTypeConverter(input string, nullable bool) string {
	dict := make(map[string]code.LangDefDataType)
	dict["varchar"] = code.LangDefDataType{Key: "varchar", NotNull: "string", Nullable: "string"}
	dict["nvarchar"] = code.LangDefDataType{Key: "varchar", NotNull: "string", Nullable: "string"}
	dict["int"] = code.LangDefDataType{Key: "varchar", NotNull: "string", Nullable: "string"}
	dict["bigint"] = code.LangDefDataType{Key: "varchar", NotNull: "string", Nullable: "string"}
	dict["bit"] = code.LangDefDataType{Key: "bool", NotNull: "string", Nullable: "string"}

	output, found := dict[input]

	if !found {
		return input
	}

	if nullable {
		return output.Nullable
	} else {
		return output.NotNull
	}

	return input
}

func snakeCase(input string) string {
	snake := wordJoiner(asWordArray(input), "_")
	return strings.ToLower(snake)
}

func kebabCase(input string) string {
	kebab := wordJoiner(asWordArray(input), "-")
	return strings.ToLower(kebab)
}

func upperCase(input string) string {
	return strings.ToUpper(input)
}

func lowerCase(input string) string {
	return strings.ToLower(input)
}

func asSentence(input string) string {
	w := asWords(input)
	o := strings.Title(w)

	return o
}
func asWords(input string) string {

	return wordJoiner(asWordArray(input), " ")
}

func pascalCase(input string) string {
	w := asWordArray(input)

	var sb strings.Builder

	for i, str := range w {

		c := strings.Title(str)
		if i == 0 {

		}
		sb.WriteString(c)
	}
	return sb.String()
}

func camelCase(input string) string {
	w := asWordArray(input)

	var sb strings.Builder

	for i, str := range w {

		c := strings.Title(str)
		if i == 0 {
			c = strings.ToLower(c)
		}
		sb.WriteString(c)
	}
	return sb.String()
}

func wordJoiner(input []string, separator string) string {
	var sb strings.Builder
	l := len(input) - 1

	for i, str := range input {
		if l == i {
			separator = ""
		}

		sb.WriteString(str + separator)
	}
	return sb.String()
}

func asWordArray(input string) []string {
	var words []string
	l := 0
	for s := input; s != ""; s = s[l:] {
		l = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
		if l <= 0 {
			l = len(s)
		}
		words = append(words, s[:l])
	}

	return words
}

func addWord(what string, input string) string {
	output := input
	if len(what) > 0 {
		output += what
	}

	return output
}