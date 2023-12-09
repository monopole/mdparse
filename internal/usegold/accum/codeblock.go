package accum

import (
	"fmt"
	"github.com/yuin/goldmark/ast"
)

type codeBlock struct {
	labels []string
	fcb    *ast.FencedCodeBlock
}

func (cb *codeBlock) printHeader(i int, content []byte) {
	fmt.Printf("%3d. %v\n", i, cb.labels)
}

func (cb *codeBlock) printCode(content []byte) {
	const (
		indent = "    "
		delim  = "```"
	)
	fmt.Print(indent)
	fmt.Print(delim)
	fmt.Println(string(cb.fcb.Language(content)))
	for i := 0; i < cb.fcb.Lines().Len(); i++ {
		s := cb.fcb.Lines().At(i)
		fmt.Print(indent)
		fmt.Print(string(s.Value(content)))
	}
	fmt.Print(indent)
	fmt.Println(delim)
}
