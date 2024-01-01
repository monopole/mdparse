package usegold

import (
	"bytes"
	"fmt"
	"github.com/monopole/mdparse/internal/loader"
	"unicode"
)

type VisitorDump2 struct {
	indent int
}

func NewVisitorDump2() *VisitorDump2 {
	return &VisitorDump2{
		indent: 0,
	}
}

const blanks = "                                                                "

func (v *VisitorDump2) VisitFolder(fl *loader.MyFolder) {
	fmt.Print(blanks[:v.indent])
	fmt.Print(fl.Name())
	if !fl.IsRoot() {
		fmt.Print("/")
	}
	fmt.Println()
	v.indent += 2
	fl.VisitFiles(v)
	fl.VisitFolders(v)
	v.indent -= 2
}

func (v *VisitorDump2) VisitFile(fi *loader.MyFile) {
	fmt.Print(blanks[:v.indent])
	fmt.Print(fi.Name())
	fmt.Print(" : ")
	fmt.Println(summarize(fi.C()) + "...")
}

func summarize(c []byte) string {
	const mx = 60
	if len(c) > mx {
		c = c[:mx]
	}
	c = bytes.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, c)
	return string(c)
}
