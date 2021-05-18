package ctx

type Context struct {
	TemplateName             string
	Templates                map[string]string
	Datatypes                map[string]string
	NullableTypes            map[string]string
	AlternativeNullableTypes map[string]string
}
