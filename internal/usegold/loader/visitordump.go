package loader

import (
	"bytes"
	"fmt"
	"path/filepath"
	"unicode"
)

type VisitorDump struct {
	indent int
}

const blanks = "                                                                "

func (v *VisitorDump) VisitContrivedFolder(f *MyContrivedFolder) {
	fmt.Printf("%s\n", f.name)
	fmt.Println("--- Origins ---")
	v.indent += 2
	fmt.Println("--- Repositories ---")
	for i := range f.repos {
		v.VisitRepo(f.repos[i])
	}
	fmt.Println("--- Absolute Tree ---")
	v.VisitFolder(f.folderAbs)
	fmt.Println("--- Process Relative Tree ---")
	v.VisitFolder(f.folderRel)
	v.indent -= 2
}

func (v *VisitorDump) VisitRepo(r *MyRepo) {
	fmt.Print(blanks[:v.indent])
	fmt.Printf("%s %s is in %s\n", r.name, r.path, r.tmpDir)
	v.indent += 2
	v.VisitFolder(r.folder)
	v.indent -= 2
}

func (v *VisitorDump) VisitFolder(fl *MyFolder) {
	fmt.Print(blanks[:v.indent])
	fmt.Print(fl.name)
	if fl.name != string(filepath.Separator) {
		fmt.Print(string(filepath.Separator))
	}
	fmt.Println()
	v.indent += 2
	for _, x := range fl.files {
		v.VisitFile(x)
	}
	for _, x := range fl.dirs {
		v.VisitFolder(x)
	}
	v.indent -= 2
}

func (v *VisitorDump) VisitFile(fi *MyFile) {
	fmt.Print(blanks[:v.indent])
	fmt.Print(fi.name)
	fmt.Print(" : ")
	c, err := fi.Contents()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(summarize(c) + "...")
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
