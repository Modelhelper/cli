package slice_test

import (
	"modelhelper/cli/slice"
	"modelhelper/cli/strstat"
	"testing"
)

func TestContains(t *testing.T) {

	list := []string{"Ola", "Kari", "Hubert", "Mike", "ruPerT"}
	cases := make(map[string]bool)
	cases["ola"] = true
	cases["karin"] = false
	cases[""] = false
	cases["ruppert"] = false

	for name, found := range cases {

		a := slice.Contains(list, name)
		e := found

		if a != e {
			t.Errorf("Contains('%s'): Expected %v, got %v", name, e, a)
		}
	}
}

func TestMaxLen(t *testing.T) {
	type testCase struct {
		expected int
		list     []string
	}
	tests := []testCase{
		{expected: 13, list: []string{"one", "two", "three", "max-len-is-13"}},
		{expected: 3, list: []string{"one", "two"}},
	}

	for idx, test := range tests {

		a := strstat.MaxLen(test.list)
		e := test.expected

		if a != e {
			t.Errorf("MaxLen of array (idx: %d): expected: %v, got %v", idx, e, a)
		}
	}
}
