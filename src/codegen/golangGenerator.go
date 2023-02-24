package codegen

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"modelhelper/cli/utils/funcmap"
	txt "modelhelper/cli/utils/text"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type GoLangGenerator struct{}

type SimpleGenerator struct{}
type simpleGenerator struct{}

// Generate implements modelhelper.TemplateGenerator
func (s *simpleGenerator) Generate(ctx context.Context, tpl *models.TextTemplate, mdl interface{}) (*models.TemplateGeneratorResult, error) {
	code, ok := ctx.Value("code").(CodeContextValue)
	res := models.TemplateGeneratorResult{}

	if !ok {
		return &res, nil
	}

	var err error
	res.Body = []byte(Generate(code.TemplateName, code.Template, mdl))
	return &res, err
}

func NewSimpleContentGenerator() modelhelper.TemplateGenerator[*models.TextTemplate] {
	return &simpleGenerator{}
}

// func (cstat *Statistics) AppendStat(instat Statistics) {
// 	cstat.Chars += instat.Chars
// 	cstat.Lines += instat.Lines
// 	cstat.Words += instat.Words
// }

func Generate(name string, body string, model interface{}) string {
	tmpl, err := template.New(name).Funcs(funcmap.SimpleFuncMap()).Parse(body)
	if err != nil {
		return ""
	}

	buf := new(bytes.Buffer)

	tmpl.Execute(buf, model)
	return buf.String()

}

func (g *GoLangGenerator) Generate(ctx context.Context, model interface{}) (*models.TemplateGeneratorResult, error) {
	start := time.Now()

	code, ok := ctx.Value("code").(CodeContextValue)
	res := models.TemplateGeneratorResult{}

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
		c, l, w := txt.GetStat(res.Body)
		res.Statistics.Chars = c
		res.Statistics.Lines = l
		res.Statistics.Words = w
		res.Statistics.Duration = time.Since(start)
	}

	return &res, nil

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

	drivers := template.Must(template.New(cv.TemplateName).Funcs(funcmap.FullFuncMap(cv.Datatypes, cv.NullableTypes)).ParseGlob(pattern))

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
