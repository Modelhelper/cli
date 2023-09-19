package tree_test

import (
	"io/ioutil"
	"modelhelper/cli/scratchpad/tree"
	"os"
	"testing"
)

func TestMaxLen(t *testing.T) {
	root := getTree()

	actual := tree.MaxLen(root)
	expected := 40

	if expected != actual {
		t.Errorf("MaxLen: expected: %d, got %d", expected, actual)
	}
}

func TestPrintTreeWithoutDesc(t *testing.T) {
	// setup test
	root := getTree()

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	tree.PrintTree(root, "", false)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	actual := string(out)
	expected := `Root Node
  Child 1
    Child 1 of Child 1
    Child 2 of Child 1
    Child 3 of Child 1 - with longest length
`

	if actual != expected {
		t.Errorf("\nExpected \n%s\n\ngot \n%s", expected, out)
	}

}
func TestPrintTreeWithDesc(t *testing.T) {
	// setup test
	root := getTreeDesc()

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	tree.PrintTree(root, "", true)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	actual := string(out)
	expected := `Root This is the root node
  C1 This is a child description
    C1.1 This is a child description
    C1.2 This is a child description
    C1.3 This is a child description
`

	if actual != expected {
		t.Errorf("\nExpected \n%s\n\ngot \n%s", expected, out)
	}

}

func getTree() tree.Node {
	root := tree.Node{Name: "Root Node"}

	c1 := tree.Node{Name: "Child 1"}
	c1.Add(tree.Node{Name: "Child 1 of Child 1"})
	c1.Add(tree.Node{Name: "Child 2 of Child 1"})
	c1.Add(tree.Node{Name: "Child 3 of Child 1 - with longest length"})
	root.Add(c1)

	return root
}
func getTreeDesc() tree.Node {
	root := tree.Node{Name: "Root", Description: "This is the root node"}

	c1 := tree.Node{Name: "C1", Description: "This is a child description"}
	c1.Add(tree.Node{Name: "C1.1", Description: "This is a child description"})
	c1.Add(tree.Node{Name: "C1.2", Description: "This is a child description"})
	c1.Add(tree.Node{Name: "C1.3", Description: "This is a child description"})
	root.Add(c1)

	return root
}
