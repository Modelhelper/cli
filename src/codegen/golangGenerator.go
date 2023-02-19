package codegen

import (
	"bufio"
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"modelhelper/cli/modelhelper"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unicode"

	"github.com/gertd/go-pluralize"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type GoLangGenerator struct{}

type SimpleGenerator struct{}

// func (cstat *Statistics) AppendStat(instat Statistics) {
// 	cstat.Chars += instat.Chars
// 	cstat.Lines += instat.Lines
// 	cstat.Words += instat.Words
// }

func Generate(name string, body string, model interface{}) string {
	tmpl, err := template.New(name).Funcs(simpleFuncMap()).Parse(body)
	if err != nil {
		return ""
	}

	buf := new(bytes.Buffer)

	tmpl.Execute(buf, model)
	return buf.String()

}

func (g *GoLangGenerator) Generate(ctx context.Context, model interface{}) (*modelhelper.CodeGeneratorResult, error) {
	start := time.Now()

	code, ok := ctx.Value("code").(CodeContextValue)
	res := modelhelper.CodeGeneratorResult{}

	if !ok {
		return &res, nil
	}
	tplMap := make(map[string]string)

	for k, b := range code.Blocks {
		tplMap[k] = b
	}
	tplMap[code.TemplateName] = code.Template

	template := fromFiles(code)
	buf := new(bytes.Buffer)
	err := template.ExecuteTemplate(buf, code.TemplateName, model)
	if err != nil {
		return nil, err
	}

	res.Body = buf.Bytes()
	if len(res.Body) > 0 {
		res.Statistics = getStat(res.Body)
		res.Statistics.Duration = time.Since(start)
	}

	return &res, nil

}

func getStat(body []byte) modelhelper.CodeGeneratorStatistics {

	s := string(body)
	stat := modelhelper.CodeGeneratorStatistics{
		Chars: len(s),
		Lines: getLines(s),
		Words: getWords(s),
	}

	return stat
}

func getWords(input string) int {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}

	return count
}
func getLines(input string) int {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanLines)

	count := 0
	for scanner.Scan() {
		count++
	}

	return count
}

func (g *SimpleGenerator) Generate(ctx context.Context, model interface{}) (*modelhelper.CodeGeneratorResult, error) {
	code, ok := ctx.Value("code").(CodeContextValue)
	res := modelhelper.CodeGeneratorResult{}

	if !ok {
		return &res, nil
	}

	var err error
	res.Body = []byte(Generate(code.TemplateName, code.Template, model))
	return &res, err

}

func fullFuncMap(dt, ntd map[string]string) template.FuncMap {
	return funcMap(stringMap(), datatypeMap(dt, ntd))
}

func simpleFuncMap() template.FuncMap {
	return funcMap(stringMap())
}

func funcMap(flist ...template.FuncMap) template.FuncMap {
	m := make(template.FuncMap)

	for _, list := range flist {
		for key, val := range list {

			m[key] = val
		}
	}

	return m
}

func datatypeMap(dt, ndt map[string]string) map[string]interface{} {
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

func stringMap() template.FuncMap {
	return template.FuncMap{
		"plural":   pluralForm,
		"singular": SingularForm,
		// "datatype":  dataTypeConverter,
		"lower":    lowerCase,
		"upper":    upperCase,
		"words":    asWords,
		"sentence": asSentence,
		"snake":    snakeCase,
		"macro":    macroCase,
		"train":    TrainCase,
		"kebab":    kebabCase,
		"dot":      DotCase,
		"title":    titleCase,
		"pascal":   pascalCase,
		"camel":    camelCase,
		// "nullable":  nullableDatatype,
		// "datatypeN": dataTypeWithNullcheck,
		"append": addWord,
	}

}

func fromFiles(cv CodeContextValue) *template.Template {
	templates := make(map[string]string)

	for k, b := range cv.Blocks {
		templates[k] = b
	}
	templates[cv.TemplateName] = cv.Template

	dir := createTempDir()
	defer os.RemoveAll(dir)

	err := writeTempFiles(dir, templates)

	if err != nil {
		return nil
	}

	pattern := filepath.Join(dir, "*")

	drivers := template.Must(template.New(cv.TemplateName).Funcs(fullFuncMap(cv.Datatypes, cv.NullableTypes)).ParseGlob(pattern))

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

func snakeCase(input string) string {
	snake := wordJoiner(asWordArray(input), "_")
	return strings.ToLower(snake)
}

func macroCase(input string) string {
	snake := wordJoiner(asWordArray(input), "_")
	return strings.ToUpper(snake)
}

func TrainCase(input string) string {
	casing := wordJoiner(asWordArray(Captial(input)), "_")
	return casing
}

func DotCase(input string) string {
	casing := wordJoiner(asWordArray(Captial(input)), ".")
	return casing
}

func kebabCase(input string) string {
	kebab := wordJoiner(asWordArray(input), "-")
	return strings.ToLower(kebab)
}
func Captial(input string) string {
	words := asWordArray(input)

	for idx, word := range words {
		word = strings.ToLower(word)
		word = strings.ToUpper(word[0:1]) + word[1:]

		words[idx] = word
	}

	return wordJoiner(words, " ")
}

func upperCase(input string) string {
	return strings.ToUpper(input)
}

func lowerCase(input string) string {
	return strings.ToLower(input)
}

func asSentence(input string) string {
	sentence := asWords(input)
	sentence = strings.ToUpper(sentence[0:1]) + strings.ToLower(sentence[1:])

	return sentence
}
func titleCase(input string) string {
	w := asWords(input)
	c := cases.Title(language.AmericanEnglish)
	return c.String(w)
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

func splitOnCasing(input string) []string {
	var words []string
	var splitPos []int
	var letterMap []int

	// nextSplitPos := 0

	// wrd := strings.Split(input, " ")
	for _, c := range input {
		val := 0
		if unicode.IsUpper(c) {
			val = 1
		}

		letterMap = append(letterMap, val)

	}

	for idx, val := range letterMap {
		if idx == 0 {
			splitPos = append(splitPos, idx)
			continue
		}

		addPos := (val == 1 && letterMap[idx-1] == 0)

		if val == 1 && idx+1 < len(letterMap) && letterMap[idx+1] == 0 {
			addPos = true
		}

		if addPos {
			splitPos = append(splitPos, idx)
		}
	}

	for idx, start := range splitPos {
		end := len(input)
		if len(splitPos) > idx+1 {
			end = splitPos[idx+1]
		}
		words = append(words, input[start:end])
	}

	return words
}

func splitOnSplitter(input string) []string {

	words := strings.FieldsFunc(input, Split)

	return words
}
func Split(r rune) bool {
	return r == ' ' || r == '_' || r == '-'
}

func asWordArray(input string) []string {
	var words []string

	split := splitOnSplitter(input)

	for _, word := range split {
		caseSplit := splitOnCasing(word)

		for _, caseWord := range caseSplit {
			words = append(words, caseWord)
		}
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
