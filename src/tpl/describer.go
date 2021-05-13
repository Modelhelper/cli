package tpl

import (
	"fmt"
	"modelhelper/cli/ui"
	"strings"
)

type Describer interface {
	Describe(key string) *Description
}

type ModelDescriber struct{}
type TypeDescriber struct{}
type GroupDescriber struct{}
type KeyDescriber struct{}

func GetDescriber(t string) Describer {
	switch strings.ToLower(t) {

	case "model":
		return &ModelDescriber{}
	case "type":
		return &TypeDescriber{}
	case "key":
		return &KeyDescriber{}
	case "group":
		return &GroupDescriber{}
	}

	return nil
}

type Description struct {
	Short string
	Long  string
}

func (d *ModelDescriber) Describe(key string) *Description {
	m, f := DefaultModels()[key]

	if f {
		return &m
	}

	return nil
}

func (d *TypeDescriber) Describe(key string) *Description {
	m, f := DefaultTypes()[key]

	if f {
		return &m
	}

	return nil
}

// func (d *GroupByType) Describe(key string) *Description {
// 	m, f := DefaultTypes()[key]

// 	if f {
// 		return &m
// 	}

// 	return nil
// }
func (d *KeyDescriber) Describe(key string) *Description {
	m, f := DefaultKeys()[key]

	if f {
		return &m
	}

	return nil
}
func (d *GroupDescriber) Describe(key string) *Description {
	m, f := DefaultModels()[key]

	if f {
		return &m
	}

	return nil
}

func (t *Template) ToString(title string) string {

	d := t.Description
	if len(d) == 0 {
		d = t.Short
	}

	o := fmt.Sprintf(`
%s

%s

Model:   %s
Key:     %s
Type:    %s
Groups:  %s
Tags:    %s

Body:
------------------------------------------------------------------------------

%s

------------------------------------------------------------------------------

Example of usage in a generate command
mh generate -t %s [options]

`,
		ui.ConsoleTitle(title),
		d,
		t.Model,
		t.Key,
		t.Type,
		strings.Join(t.Groups, ", "),
		strings.Join(t.Tags, ", "),
		t.Body,
		title,
	)

	return o
}
