package accum

import (
	"bytes"
	"fmt"
	"github.com/monopole/mdrip/base"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"os"
	"strings"
)

// LessonDoc groups a markdown AST with the codeblocks found within in.
// LessonDoc has a one to one correspondence to a file.
// It must have a name, and may have blocks.
// An entirely empty file might appear with no blocks.
type LessonDoc struct {
	// Where was the content read from if it was read from a path?
	// This might be empty, if the content came from, say, a stream.
	path base.FilePath
	// The raw content read from the path.
	content []byte
	// An abstract syntax tree discovered by parsing the content.
	// Cannot be used alone, as it holds pointers into content.
	doc ast.Node
	// The code blocks found in the AST.
	blocks []*CodeBlock
}

func NewLessonDocFromPath(p parser.Parser, path base.FilePath) (ld *LessonDoc, err error) {
	var c []byte
	c, err = os.ReadFile(string(path))
	if err != nil {
		return
	}
	ld, err = NewLessonDocFromBytes(p, c)
	ld.path = path
	return
}

func NewLessonDocFromBytes(p parser.Parser, c []byte) (ld *LessonDoc, err error) {
	ld = &LessonDoc{
		content: c,
		doc:     p.Parse(text.NewReader(c)),
	}
	err = ld.RefreshBlocks()
	return
}

func (ld *LessonDoc) Dump() {
	ld.doc.Dump(ld.content, 0)
	ld.printBlocks()
}

func (ld *LessonDoc) Render(r renderer.Renderer) (string, error) {
	var b bytes.Buffer
	err := r.Render(&b, ld.content, ld.doc)
	return b.String(), err
}

func (ld *LessonDoc) RefreshBlocks() error {
	ld.blocks = nil
	return ast.Walk(ld.doc, ld.walkForBlocks)
}

func (ld *LessonDoc) walkForBlocks(n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	if n.Kind() == ast.KindFencedCodeBlock {
		if fcb, ok := n.(*ast.FencedCodeBlock); ok {
			ld.accumulateCodeBlock(fcb)
		} else {
			return ast.WalkStop, fmt.Errorf("ast.Kind() is dishonest")
		}
	}
	return ast.WalkContinue, nil
}

func (ld *LessonDoc) accumulateCodeBlock(fcb *ast.FencedCodeBlock) {
	cb := NewCodeBlock(ld, fcb)
	if prev := fcb.PreviousSibling(); prev != nil && prev.Kind() == ast.KindHTMLBlock {
		if html, ok := prev.(*ast.HTMLBlock); ok {
			// We have a preceding HTML block.
			// If it's an HTML comment, try to extract labels.
			cb.labels = append(cb.labels, parseLabels(commentBody(ld.nodeText(html)))...)
		}
	}
	ld.blocks = append(ld.blocks, cb)
}

func (ld *LessonDoc) printBlocks() {
	for i, b := range ld.blocks {
		b.printHeader(i, ld.content)
		b.printCode()
	}
}

func (ld *LessonDoc) nodeText(n ast.Node) string {
	var buff strings.Builder
	for i := 0; i < n.Lines().Len(); i++ {
		s := n.Lines().At(i)
		buff.Write(ld.content[s.Start:s.Stop])
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

func parseLabels(s string) (result []base.Label) {
	const labelPrefixChar = uint8('@')
	items := strings.Split(s, " ")
	for _, word := range items {
		i := 0
		for i < len(word) && word[i] == labelPrefixChar {
			i++
		}
		if i > 0 && i < len(word) && word[i-1] == labelPrefixChar {
			result = append(result, base.Label(word[i:]))
		}
	}
	return
}

// Accept accepts a visitor.
func (ld *LessonDoc) Accept(v TutVisitor) { v.VisitLessonDoc(ld) }

// Title is the purported title of the LessonDoc.
func (ld *LessonDoc) Title() string {
	// TODO try to grab this from the AST.
	return ld.Name()
}

// Name is the purported name of the lesson.
func (ld *LessonDoc) Name() string {
	if ld.Path().IsEmpty() {
		return "NO IDEA WHAT THIS IS CALLED"
	}
	return ld.path.Base()
}

// Path to the lesson.  A lesson has a 1:1 correspondence with a path.
func (ld *LessonDoc) Path() base.FilePath { return ld.path }

// Children of the lesson - the code blocks.
func (ld *LessonDoc) Children() (result []Tutorial) {
	for _, b := range ld.blocks {
		result = append(result, b)
	}
	return result
}
