package accum

import (
	"fmt"
	"github.com/monopole/mdrip/util"
	"io"
)

// TxtPrinter prints a tutorial as text.
type TxtPrinter struct {
	depth int
	w     io.Writer
}

// NewTutorialTxtPrinter makes a new TxtPrinter for the given writer.
func NewTutorialTxtPrinter(w io.Writer) *TxtPrinter {
	return &TxtPrinter{0, w}
}

func (v *TxtPrinter) wrapFmt(s string) string {
	return util.Spaces(2*v.depth) + s + "\n"
}

// Depth is how deep we are in a tutorial tree.
func (v *TxtPrinter) Depth() int {
	return v.depth
}

// Down goes deeper.
func (v *TxtPrinter) Down() {
	v.depth++
}

// Up is opposite of Down.
func (v *TxtPrinter) Up() {
	v.depth--
}

// P does a formatted print.
func (v *TxtPrinter) P(s string, a ...interface{}) {
	fmt.Fprintf(v.w, v.wrapFmt(s), a...)
}

// VisitCodeBlock prints a CodeBlock.
func (v *TxtPrinter) VisitCodeBlock(b *CodeBlock) {
	v.P("%s --- %s...", b.Name(), util.SampleString(string(b.Code()), 60))
}

// VisitLessonDoc prints a LessonDoc.
func (v *TxtPrinter) VisitLessonDoc(l *LessonDoc) {
	v.P("%s", l.Name())
	v.Down()
	for _, x := range l.Children() {
		x.Accept(v)
	}
	v.Up()
}

// VisitCourse prints a Course.
func (v *TxtPrinter) VisitCourse(c *Course) {
	v.P("%s", c.Name())
	v.Down()
	for _, x := range c.Children() {
		x.Accept(v)
	}
	v.Up()
}
