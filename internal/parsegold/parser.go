package parsegold

import (
	bytes2 "bytes"
	"fmt"
	"github.com/monopole/mdparse/internal/ifc"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

type gomark struct {
	p    goldmark.Markdown
	data []byte
	doc  ast.Node
}

func NewMarkdownParser(doMyStuff bool) ifc.Marker {
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	return &gomark{p: markdown}

}

func (gm *gomark) Parse(bytes []byte) error {
	gm.data = bytes

	var doc ast.Node

	doc = gm.p.Parser().Parse(text.NewReader(bytes))
	// doc.Meta()["footnote-prefix"] = getPrefix(path)
	fmt.Printf("%T %+v\n", doc, doc)
	doc.Dump(bytes, 0)
	// Dump and Render need the original source text because the AST doesn't
	// hold the original text - it just has byte array offsets.
	// Every Node is a BaseBlock, and each BaseBlock has a ptr to the lines in the
	// source text that make it, and each line is a Segment, and each Segment has
	// a Start and Stop integer index meant for use with a byte array.
	// I confirmed this by sending some different document in.

	return nil
}

func (gm *gomark) Render() (string, error) {
	var b bytes2.Buffer
	err := gm.p.Renderer().Render(&b, gm.data, gm.doc)
	return b.String(), err
}
