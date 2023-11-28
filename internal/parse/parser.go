package parse

import (
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

func NewMarkdownParser(doMyStuff bool) *parser.Parser {
	p := parser.NewWithExtensions(parser.CommonExtensions |
		parser.AutoHeadingIDs |
		parser.NoEmptyLineBeforeBlock |
		parser.Attributes)
	if doMyStuff {
		p.Opts.ParserHook = parserHook
	}
	return p
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
