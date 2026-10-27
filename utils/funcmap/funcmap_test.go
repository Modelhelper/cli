package funcmap

import (
	"reflect"
	"testing"

	"github.com/Modelhelper/cli/v3/modelhelper/models"
)

func TestLanguageMap(t *testing.T) {
	lang := &models.LanguageDefinition{
		ReservedWords: []string{"type"},
		DataTypes: map[string]models.Datatype{
			"datetime": {Imports: []string{"time"}, NullableImports: []string{"database/sql"}},
			"varchar":  {NullableImports: []string{"database/sql"}},
		},
	}
	fm := LanguageMap(lang)
	imports := fm["imports"].(func([]models.EntityColumnModel) []string)
	safe := fm["safe"].(func(string) string)

	cols := []models.EntityColumnModel{
		{DataType: "datetime"},
		{DataType: "datetime"},
		{DataType: "varchar", IsNullable: true},
		{DataType: "varchar"},
		{DataType: "unknown"},
	}

	if got, want := imports(cols), []string{"database/sql", "time"}; !reflect.DeepEqual(got, want) {
		t.Errorf("imports = %v, want %v", got, want)
	}

	if got := safe("type"); got != "type_" {
		t.Errorf("safe(type) = %s", got)
	}
	if got := safe("name"); got != "name" {
		t.Errorf("safe(name) = %s", got)
	}

	// no language: nothing to import, names are kept
	none := LanguageMap(nil)
	if got := none["imports"].(func([]models.EntityColumnModel) []string)(cols); len(got) != 0 {
		t.Errorf("imports without a language = %v", got)
	}
	if got := none["safe"].(func(string) string)("type"); got != "type" {
		t.Errorf("safe without a language = %s", got)
	}
}
