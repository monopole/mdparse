package main

import (
	"fmt"
	"github.com/gomarkdown/markdown/ast"
	"strings"
)

func myWalk(doc ast.Node) {
	fmt.Println("Walking...")
	ast.Walk(doc, &nodeVisitor{})
	fmt.Println("Done Walking.")
}

type nodeVisitor struct {
	indent string
}

func (v *nodeVisitor) Visit(n ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.GoToNext
	}
	// ast.Print recurses its argument, instead of just visiting
	// only the argument, so it's not what you want.
	// ast.Print(os.Stdout, n)
	fmt.Println(nodeType(n))
	return ast.GoToNext
}

// get a short name of the type of v which excludes package name
// and strips "()" from the end
func nodeType(node ast.Node) string {
	s := fmt.Sprintf("%T", node)
	s = strings.TrimSuffix(s, "()")
	if idx := strings.Index(s, "."); idx != -1 {
		return s[idx+1:]
	}
	return s
}
