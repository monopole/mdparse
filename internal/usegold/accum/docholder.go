package accum

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"strings"
)

// docHolder associates an AST tree with the raw content that
// was parsed to form the tree.
// The AST doesn't hold the original text - it just has byte array offsets.
// Every ast.Node is a BaseBlock, and each BaseBlock has a ptr to the lines in the
// source text that make it, and each line is a Segment, and each Segment has
// a Start and Stop integer index meant for use with a byte array.
type docHolder struct {
	doc     ast.Node
	content []byte
	blocks  []*codeBlock
}

func (dh *docHolder) Dump() {
	dh.doc.Dump(dh.content, 0)
	dh.printBlocks()
}

func (dh *docHolder) Render(r renderer.Renderer) (string, error) {
	var b bytes.Buffer
	err := r.Render(&b, dh.content, dh.doc)
	return b.String(), err
}

func (dh *docHolder) myWalk(n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	if n.Kind() == ast.KindFencedCodeBlock {
		if fcb, ok := n.(*ast.FencedCodeBlock); ok {
			dh.accumulateCodeBlock(fcb)
		} else {
			return ast.WalkStop, fmt.Errorf("ast.Kind() is dishonest")
		}
	}
	return ast.WalkContinue, nil
}

func (dh *docHolder) accumulateCodeBlock(fcb *ast.FencedCodeBlock) {
	cb := &codeBlock{fcb: fcb}
	if prev := fcb.PreviousSibling(); prev != nil && prev.Kind() == ast.KindHTMLBlock {
		if html, ok := prev.(*ast.HTMLBlock); ok {
			// We have a preceding HTML block.
			// If it's an HTML comment, try to extract labels.
			cb.labels = append(cb.labels, parseLabels(commentBody(dh.nodeText(html)))...)
		}
	}
	dh.blocks = append(dh.blocks, cb)
}

func (dh *docHolder) printBlocks() {
	for i, b := range dh.blocks {
		b.printHeader(i, dh.content)
		b.printCode(dh.content)
	}
}

func (dh *docHolder) nodeText(n ast.Node) string {
	var buff strings.Builder
	for i := 0; i < n.Lines().Len(); i++ {
		s := n.Lines().At(i)
		buff.Write(dh.content[s.Start:s.Stop])
	}
	return buff.String()
}

func commentBody(s string) string {
	const (
		begin = "<!--"
		end   = "-->"
	)
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, begin) {
		return ""
	}
	if !strings.HasSuffix(s, end) {
		return ""
	}
	return s[len(begin) : len(s)-len(end)]
}

func parseLabels(s string) (result []string) {
	const labelPrefixChar = uint8('@')
	items := strings.Split(s, " ")
	for _, word := range items {
		i := 0
		for i < len(word) && word[i] == labelPrefixChar {
			i++
		}
		if i > 0 && i < len(word) && word[i-1] == labelPrefixChar {
			result = append(result, word[i:])
		}
	}
	return
}
