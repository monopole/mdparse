package loader

import (
	"bytes"
	"fmt"
	"unicode"
)

type VisitorDump struct {
	indent int
}

func NewVisitorDump() *VisitorDump {
	return &VisitorDump{
		indent: 0,
	}
}

const blanks = "                                                                "

func (v *VisitorDump) VisitFolder(fl *MyFolder) {
	fmt.Print(blanks[:v.indent])
	fmt.Print(fl.Name())
	if !fl.IsRoot() {
		fmt.Print(rootSlash)
	}
	fmt.Println()
	v.indent += 2
	fl.VisitFiles(v)
	fl.VisitFolders(v)
	v.indent -= 2
}

func (v *VisitorDump) VisitFile(fi *MyFile) {
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
