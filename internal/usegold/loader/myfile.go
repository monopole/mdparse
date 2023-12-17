package loader

import (
	"os"
)

// MyFile is named byte array.
type MyFile struct {
	myTreeItem
	content []byte
}

var _ MyTreeItem = &MyFile{}

func (fi *MyFile) Accept(v TreeVisitor) {
	v.VisitFile(fi)
}

//// myErrFile returns "fake" markdown file showing an error message.
//func myErrFile(path string, err error) (*MyFile, error) {
//	return &MyFile{
//		myTreeItem: myTreeItem{
//			parent: nil,
//			name:   path,
//		},
//		content: []byte(fmt.Sprintf("## Unable to load from %s; %s", path, err.Error())),
//	}, err
//}

func (fi *MyFile) Contents() ([]byte, error) {
	return os.ReadFile(fi.FullName())
}

func ReorderFiles(x []*MyFile, ordering []string) []*MyFile {
	for i := len(ordering) - 1; i >= 0; i-- {
		x = ShiftFileToTop(x, ordering[i])
	}
	return ShiftFileToTop(x, "README")
}

func ShiftFileToTop(x []*MyFile, top string) []*MyFile {
	var first []*MyFile
	var remainder []*MyFile
	for _, f := range x {
		if f.Name() == top {
			first = append(first, f)
		} else {
			remainder = append(remainder, f)
		}
	}
	return append(first, remainder...)
}
