package accum

import (
	"fmt"
	"github.com/monopole/mdrip/base"
	"github.com/yuin/goldmark/ast"
	"strings"
)

// CodeBlock groups an ast.FencedCodeBlock with a set of labels.
type CodeBlock struct {
	labels []base.Label
	fcb    *ast.FencedCodeBlock
	parent *LessonDoc
}

func NewCodeBlock(ld *LessonDoc, fcb *ast.FencedCodeBlock) *CodeBlock {
	return &CodeBlock{fcb: fcb, parent: ld}
}

func (cb *CodeBlock) printHeader(i int, content []byte) {
	fmt.Printf("%3d. %v\n", i, cb.labels)
}

func (cb *CodeBlock) printCode() {
	const (
		indent = "    "
		delim  = "```"
	)
	fmt.Print(indent)
	fmt.Print(delim)
	fmt.Println(string(cb.fcb.Language(cb.parent.content)))
	for i := 0; i < cb.fcb.Lines().Len(); i++ {
		s := cb.fcb.Lines().At(i)
		fmt.Print(indent)
		fmt.Print(string(s.Value(cb.parent.content)))
	}
	fmt.Print(indent)
	fmt.Println(delim)
}

func (cb *CodeBlock) Code() string {
	var b strings.Builder
	for i := 0; i < cb.fcb.Lines().Len(); i++ {
		s := cb.fcb.Lines().At(i)
		b.Write(s.Value(cb.parent.content))
	}
	return b.String()
}

// HasLabel is true if the block has the given label argument.
func (cb *CodeBlock) HasLabel(label base.Label) bool {
	for _, l := range cb.labels {
		if l == label {
			return true
		}
	}
	return false
}

// AnonBlockName used for blocks that have no explicit name.
const AnonBlockName = "clickToCopy"

// Title is what appears to be the title of the block.
func (cb *CodeBlock) Title() string {
	return cb.Name()
}

// Name attempts to return a decent name for the block.
func (cb *CodeBlock) Name() string {
	l := cb.firstNiceLabel()
	if l == base.AnonLabel {
		return AnonBlockName
	}
	return string(l)
}

func (cb *CodeBlock) firstNiceLabel() base.Label {
	for _, l := range cb.labels {
		if l != base.WildCardLabel && l != base.AnonLabel {
			return l
		}
	}
	return base.AnonLabel
}

// Accept accepts a visitor.
func (cb *CodeBlock) Accept(v TutVisitor) { v.VisitCodeBlock(cb) }

// Path to the file containing the block.
func (cb *CodeBlock) Path() base.FilePath { return base.FilePath("notUsingThis") }

// Children of the block - there aren't any at this time.
// One could imagine each line of code in a code block as a child
// if that were useful somehow.
func (cb *CodeBlock) Children() []Tutorial { return nil }
