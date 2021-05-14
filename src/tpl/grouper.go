package tpl

import "strings"

type TemplateGrouper interface {
	Group(t map[string]Template) map[string]map[string]Template
}

func GetGrouper(t string) TemplateGrouper {

	switch strings.ToLower(t) {

	case "language":
		return &GroupByLanguage{}
	case "group":
		return &GroupByGroup{}
	case "key":
		return &GroupByKey{}
	case "model":
		return &GroupByModel{}
	case "tag":
		return &GroupByTag{}
	default:
		return &GroupByType{}

	}
}

type GroupByType struct{}
type GroupByGroup struct{}
type GroupByTag struct{}
type GroupByLanguage struct{}
type GroupByKey struct{}
type GroupByModel struct{}

func (g *GroupByType) Group(t map[string]Template) map[string]map[string]Template {
	m := make(map[string]map[string]Template)

	for tn, template := range t {

		k, f := m[template.Type]

		if !f {
			k = make(map[string]Template)
		}

		k[tn] = template
		m[template.Type] = k

	}

	return m
}
func (g *GroupByKey) Group(t map[string]Template) map[string]map[string]Template {
	m := make(map[string]map[string]Template)
	empty := make(map[string]Template)

	for tn, template := range t {

		if len(template.Key) == 0 {
			empty[tn] = template
		} else {

			k, f := m[template.Key]

			if !f {
				k = make(map[string]Template)
			}

			k[tn] = template
			m[template.Key] = k
		}

	}

	if len(empty) > 0 {
		m["nokey"] = empty
	}

	return m
}
func (g *GroupByModel) Group(t map[string]Template) map[string]map[string]Template {
	m := make(map[string]map[string]Template)
	empty := make(map[string]Template)

	for tn, template := range t {

		if len(template.Model) == 0 {
			empty[tn] = template
		} else {

			k, f := m[template.Model]

			if !f {
				k = make(map[string]Template)
			}

			k[tn] = template
			m[template.Model] = k
		}

	}

	if len(empty) > 0 {
		m["empty"] = empty
	}

	return m
}

func (g *GroupByLanguage) Group(t map[string]Template) map[string]map[string]Template {
	m := make(map[string]map[string]Template)
	empty := make(map[string]Template)

	for tn, template := range t {

		if len(template.Language) > 0 {
			empty[tn] = template
		} else {

			if template.Type != "block" {

				k, f := m[template.Language]

				if !f {
					k = make(map[string]Template)
				}

				k[tn] = template
				m[template.Language] = k
			}
		}
	}

	if len(empty) > 0 {
		m["empty"] = empty
	}

	return m
}
func (g *GroupByTag) Group(t map[string]Template) map[string]map[string]Template {
	m := make(map[string]map[string]Template)
	empty := make(map[string]Template)

	for tn, template := range t {

		if template.Type != "block" {

			var tmpl = template

			if len(template.Tags) == 0 {
				empty[tn] = tmpl
			} else {
				for _, grp := range template.Tags {

					k, f := m[grp]

					if !f {
						k = make(map[string]Template)
					}

					k[tn] = tmpl
					m[grp] = k
				}
			}
		}

	}

	if len(empty) > 0 {
		m["empty"] = empty
	}
	return m
}
func (g *GroupByGroup) Group(t map[string]Template) map[string]map[string]Template {
	m := make(map[string]map[string]Template)
	empty := make(map[string]Template)

	for tn, template := range t {

		if template.Type != "block" {

			var tmpl = template

			if len(template.Groups) == 0 {
				empty[tn] = tmpl
			} else {
				for _, grp := range template.Groups {

					k, f := m[grp]

					if !f {
						k = make(map[string]Template)
					}

					k[tn] = tmpl
					m[grp] = k
				}
			}
		}
	}

	if len(empty) > 0 {
		m["empty"] = empty
	}

	return m
}
