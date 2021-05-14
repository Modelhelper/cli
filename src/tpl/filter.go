package tpl

type TemplateFilter interface {
	Filter(t map[string]Template, filter []string) map[string]Template
}
type FilterByLang struct{}
type FilterByType struct{}
type FilterByModel struct{}
type FilterByKey struct{}
type FilterByGroup struct{}

func (f *FilterByType) Filter(t map[string]Template, filter []string) map[string]Template {
	output := make(map[string]Template)

	for name, template := range t {
		if contains(filter, template.Type) {
			output[name] = template
		}
	}

	return output
}
func (f *FilterByLang) Filter(t map[string]Template, filter []string) map[string]Template {
	output := make(map[string]Template)

	for name, template := range t {
		if contains(filter, template.Language) {
			output[name] = template
		}
	}

	return output
}
func (f *FilterByModel) Filter(t map[string]Template, filter []string) map[string]Template {
	output := make(map[string]Template)

	for name, template := range t {
		if contains(filter, template.Model) {
			output[name] = template
		}
	}

	return output
}
func (f *FilterByGroup) Filter(t map[string]Template, filter []string) map[string]Template {
	output := make(map[string]Template)

	for name, template := range t {
		if len(template.Groups) > 0 {
			for _, grp := range template.Groups {

				if contains(filter, grp) {
					output[name] = template
				}
			}
		}
	}

	return output
}
func (f *FilterByKey) Filter(t map[string]Template, filter []string) map[string]Template {
	output := make(map[string]Template)

	for name, template := range t {
		if contains(filter, template.Key) {
			output[name] = template
		}
	}

	return output
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
