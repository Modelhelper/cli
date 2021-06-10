package codegen

import (
	"context"
	"modelhelper/cli/model"
	"modelhelper/cli/project"
	"time"
)

// var datatypeMap map[string]code.LangDefDataType

// type CodeGenerator interface {
// 	Generate(model interface{}) (string, error)
// }

type CodeGen interface {
	Generate(ctx context.Context, m model.ModelConverter) (Result, error)
}

type CodeContextValue struct {
	TemplateName  string
	Template      string
	Blocks        map[string]string
	Datatypes     map[string]string
	NullableTypes map[string]string
	Project       project.Project

	// AlternativeNullableTypes map[string]string
	// Templates                map[string]string // is this really needed
}

type Statistics struct {
	Chars     int
	Lines     int
	Words     int
	Duration  time.Duration
	TimeSaved int
}

type Result struct {
	Stat    Statistics
	Content string
}
