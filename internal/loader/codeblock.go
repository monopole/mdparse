package loader

import (
	"fmt"
	"github.com/monopole/mdrip/base"
)

// CodeBlock groups an ast.FencedCodeBlock with a set of labels.
type CodeBlock struct {
	labels   []base.Label
	code     string
	language string
	parent   *MyFile
}

func NewCodeBlock(
	fi *MyFile, code, language string) *CodeBlock {
	return &CodeBlock{code: code, language: language, parent: fi}
}

func (cb *CodeBlock) printHeader(i int, content []byte) {
	fmt.Printf("%3d. %v\n", i, cb.labels)
}

func (cb *CodeBlock) AddLabels(labels []base.Label) {
	cb.labels = append(cb.labels, labels...)
}

func (cb *CodeBlock) Code() string {
	return cb.code
}

func (cb *CodeBlock) Dump() {
	if len(cb.labels) > 0 {
		fmt.Print("# labels: ")
		for _, l := range cb.labels {
			fmt.Print(" ", l)
		}
		fmt.Println()
	}
	fmt.Printf("# lang=%q\n", cb.language)
	fmt.Print(cb.code)
	fmt.Println("# -----------")
}

// HasLabel is true if the block has the given label argument.
func (cb *CodeBlock) HasLabel(label base.Label) bool {
	if label == base.WildCardLabel {
		return true
	}
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
