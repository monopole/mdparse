package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

type Gallery struct {
	ast.Leaf
	ImageURLS []string
}

var gallery = []byte(":gallery\n")

func parseGallery(data []byte) (ast.Node, []byte, int) {
	if !bytes.HasPrefix(data, gallery) {
		return nil, nil, 0
	}
	fmt.Printf("Found a gallery!\n\n")
	i := len(gallery)
	// find empty line
	// TODO: should also consider end of document
	end := bytes.Index(data[i:], []byte("\n\n"))
	if end < 0 {
		return nil, data, 0
	}
	end = end + i
	lines := string(data[i:end])
	parts := strings.Split(lines, "\n")
	res := &Gallery{
		ImageURLS: parts,
	}
	return res, nil, end
}

func parserHook(data []byte) (ast.Node, []byte, int) {
	if node, d, n := parseGallery(data); node != nil {
		return node, d, n
	}
	return nil, nil, 0
}

const doMyStuff = false

func newMarkdownParser() *parser.Parser {
	extensions := parser.CommonExtensions |
		parser.AutoHeadingIDs |
		parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	if doMyStuff {
		p.Opts.ParserHook = parserHook
	}
	return p
}
