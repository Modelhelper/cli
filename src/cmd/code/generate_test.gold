package cmd

import "testing"

func TestCoalesce(t *testing.T) {
	actual := coalesceString("", "org")
	expected := "org"

	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

	actual = coalesceString("this", "org")
	expected = "this"

	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
