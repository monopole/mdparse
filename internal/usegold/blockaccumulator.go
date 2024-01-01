package usegold

import (
	"fmt"
	"github.com/monopole/mdparse/internal/loader"
	"github.com/monopole/mdrip/base"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"strings"
)

// BlockAccumulator uses the goldmark parser to find code blocks.
type BlockAccumulator struct {
	p           goldmark.Markdown
	currentFile *loader.MyFile

	// The code blocks found in the AST.
	blocks []*loader.CodeBlock
}

func NewBlockAccumulator() *BlockAccumulator {
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
	return &BlockAccumulator{
		p: markdown,
	}
}

const blanks = "                                                                "

func (v *BlockAccumulator) Blocks(l base.Label) []*loader.CodeBlock {
	var result []*loader.CodeBlock
	for i := range v.blocks {
		if v.blocks[i].HasLabel(l) {
			result = append(result, v.blocks[i])
		}
	}
	return result
}

func (v *BlockAccumulator) VisitFolder(fl *loader.MyFolder) {
	fl.VisitFiles(v)
	fl.VisitFolders(v)
}

func (v *BlockAccumulator) VisitFile(fi *loader.MyFile) {
	v.currentFile = fi
	// An abstract syntax tree discovered by parsing the content.
	// Cannot be used alone, as it holds pointers into content.
	doc := v.p.Parser().Parse(text.NewReader(fi.C()))
	fmt.Printf("scanning %s\n", fi.Name())
	ast.Walk(doc, v.walkForBlocks)
}

func (v *BlockAccumulator) walkForBlocks(n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	if n.Kind() == ast.KindFencedCodeBlock {
		if fcb, ok := n.(*ast.FencedCodeBlock); ok {
			v.accumulateCodeBlock(fcb)
		} else {
			return ast.WalkStop, fmt.Errorf("ast.Kind() is dishonest")
		}
	}
	return ast.WalkContinue, nil
}

func (v *BlockAccumulator) accumulateCodeBlock(fcb *ast.FencedCodeBlock) {
	cb := loader.NewCodeBlock(
		v.currentFile, v.nodeText(fcb), string(fcb.Language(v.currentFile.C())))
	if prev := fcb.PreviousSibling(); prev != nil && prev.Kind() == ast.KindHTMLBlock {
		if html, ok := prev.(*ast.HTMLBlock); ok {
			// We have a preceding HTML block.
			// If it's an HTML comment, try to extract labels.
			cb.AddLabels(loader.ParseLabels(loader.CommentBody(v.nodeText(html))))
		}
	}
	v.blocks = append(v.blocks, cb)
}

// TODO: Could change this to preserve lines?
func (v *BlockAccumulator) nodeText(n ast.Node) string {
	var buff strings.Builder
	for i := 0; i < n.Lines().Len(); i++ {
		s := n.Lines().At(i)
		buff.Write(v.currentFile.C()[s.Start:s.Stop])
	}
	return buff.String()
}
