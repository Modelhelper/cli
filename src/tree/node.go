//tree package works with nodes that can
package tree

import (
	"fmt"
)

type TreeBuilder interface {
	Build() Node
}

type Node struct {
	// ID          int
	Name        string
	Description string
	Nodes       []Node
	// RelColumn   string
	// RelName     string
}

func (n *Node) Add(child Node) {
	n.Nodes = append(n.Nodes, child)
}

func PrintTree(root Node, prefix string, printDescription bool) {
	if prefix == "" {
		if printDescription && len(root.Description) > 0 {
			fmt.Printf("%s %s\n", root.Name, root.Description)
		} else {
			fmt.Printf("%s\n", root.Name)
		}
	} else {
		if printDescription {
			fmt.Printf("%s%s %s\n", prefix, root.Name, root.Description)
		} else {
			fmt.Printf("%s%s\n", prefix, root.Name)
		}
	}
	for _, n := range root.Nodes {
		PrintTree(n, prefix+"  ", printDescription)
	}
}

func MaxLen(root Node) int {

	l := len(root.Name)
	for _, n := range root.Nodes {
		// tl := len(n.Name)
		tl := MaxLen(n)
		if l < tl {
			l = tl
		}
	}

	return l
}
