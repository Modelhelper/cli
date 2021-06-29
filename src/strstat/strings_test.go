package strstat_test

import (
	"modelhelper/cli/strstat"
	"testing"
)

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

// func BenchMaxLen(t *testing.B) {

// }
