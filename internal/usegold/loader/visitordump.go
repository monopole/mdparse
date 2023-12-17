package loader

import (
	"fmt"
	"path/filepath"
)

type VisitorDump struct {
	indent int
}

const blanks = "                                                                "

func (v *VisitorDump) VisitContrivedFolder(f *MyContrivedFolder) {
	fmt.Printf("%s, nRepos=%d\n", f.name, len(f.repos))
	for i := range f.originalSpecs {
		fmt.Println("  ", f.originalSpecs[i])
	}
	for i := range f.repos {
		v.VisitRepo(f.repos[i])
	}
	fmt.Println("Absolute Tree")
	v.VisitFolder(f.folderAbs)
	fmt.Println("Relative Tree")
	v.VisitFolder(f.folderRel)
}

func (v *VisitorDump) VisitRepo(_ *MyRepo) {
	fmt.Print("TODO: dump repo")
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
	fmt.Println(fi.name)
}
