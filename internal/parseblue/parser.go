package parseblue

import (
	"fmt"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/monopole/mdparse/internal/ifc"
	"github.com/monopole/mdparse/internal/render"
	"os"
	"strings"
)

type gomark struct {
	doMyStuff bool
	p         *parser.Parser
	doc       ast.Node
}

func (gm *gomark) Parse(data []byte) error {
	gm.doc = gm.p.Parse(data)

	ast.PrintWithPrefix(os.Stdout, gm.doc, "  ")
	myWalk(gm.doc)
	//	_, err = fmt.Printf("--- Markdown:\n%s\n\n", md)
	return nil
}

func (gm *gomark) Render() (string, error) {
	res := render.RenderAsHtml(gm.doc, gm.doMyStuff)
	return string(res), nil
}

func NewMarkdownParser(doMyStuff bool) ifc.Marker {
	p := parser.NewWithExtensions(parser.CommonExtensions |
		parser.AutoHeadingIDs |
		parser.NoEmptyLineBeforeBlock |
		parser.Attributes)
	if doMyStuff {
		p.Opts.ParserHook = parserHook
	}
	return &gomark{p: p}
}

// parserHook is a custom parser.
// If successful it returns an ast.Node containing the results of the parsing,
// a buffer that should be parsed as a block and added to the document (see below),
// and the number of bytes consumed (the guts of the parser will skip over this).
// The buffer returned could be anything - e.g. data pulled from the web.
// Any nodes parsed from it will follow
// the ast.Node returned here at the same level, and not be a child to it.
// I think this is normally nil?
// It seems to be a way to inject data into the document.
func parserHook(data []byte) (ast.Node, []byte, int) {
	if node, d, n := attemptToParseGallery(data); node != nil {
		return node, d, n
	}
	return nil, nil, 0
}

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
	leafLiteral := ""
	if n.AsLeaf() != nil {
		leafLiteral = string(n.AsLeaf().Literal)
	}
	fmt.Printf("%s %s \n", nodeType(n), leafLiteral)
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
