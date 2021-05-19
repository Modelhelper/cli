package codegen

import (
	"bytes"
	"io/ioutil"
	"log"
	"modelhelper/cli/ctx"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/gertd/go-pluralize"
)

var codeCtx ctx.Context

type GoLangGenerator struct {
	TemplateName  string
	Templates     map[string]string
	Datatypes     map[string]string
	NullableTypes map[string]string
}

type SimpleGenerator struct {
	Template string
}

func Generate(name string, body string, model interface{}) (string, error) {
	tmpl, err := template.New(name).Funcs(funcMap()).Parse(body)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	tmpl.Execute(buf, model)
	return buf.String(), nil

}

func (g *GoLangGenerator) Generate(c ctx.Context, model interface{}) (string, error) {
	codeCtx = c

	template := fromFiles(g.TemplateName, g.Templates)

	buf := new(bytes.Buffer)
	err := template.ExecuteTemplate(buf, g.TemplateName, model)

	if err != nil {
		// fmt.Println(err)
		log.Fatalln(err)
	}

	return buf.String(), nil

}
func (g *SimpleGenerator) Generate(c ctx.Context, model interface{}) (string, error) {
	codeCtx = c

	if len(c.Templates) > 0 {
		return Generate(c.TemplateName, g.Template, model)
	}

	return "", nil
}

func funcMap() template.FuncMap {
	return template.FuncMap{
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

}

func fromFiles(name string, templates map[string]string) *template.Template {

	dir := createTempDir()
	defer os.RemoveAll(dir)

	err := writeTempFiles(dir, templates)

	if err != nil {
		return nil
	}

	pattern := filepath.Join(dir, "*")

	drivers := template.Must(template.New(name).Funcs(funcMap()).ParseGlob(pattern))

	return drivers
}

func createTempDir() string {

	dir, err := ioutil.TempDir("", "template") //os.MkdirTemp("", "template")
	if err != nil {
		log.Fatal(err)
	}

	return dir

}

func writeTempFiles(dir string, templates map[string]string) error {
	for name, body := range templates {
		fp := filepath.Join(dir, name)
		err := ioutil.WriteFile(fp, []byte(body), 0777)
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

	dt := dataTypeConverter(input)

	// if reflect.TypeOf(isNullable) == reflect.Typeof(bool)
	if isNullable {
		return nullableDatatype(dt)
	}
	return dt
}
func alternativeNullableDatatype(input string) string {
	output := codeCtx.AlternativeNullableTypes[input]

	if len(output) > 0 {
		return output
	}

	return input
}
func nullableDatatype(input string) string {
	output, found := codeCtx.NullableTypes[input]

	if found {
		return output
	}

	return input
}

func dataTypeConverter(input string) string {

	output, found := codeCtx.Datatypes[input]

	if found {
		return output
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
