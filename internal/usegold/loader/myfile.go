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

func NewFile(n string) *MyFile {
	return &MyFile{myTreeItem: myTreeItem{name: n}}
}

func (fi *MyFile) Accept(v TreeVisitor) {
	v.VisitFile(fi)
}

func (fi *MyFile) Equals(other *MyFile) bool {
	if fi == nil {
		return other == nil
	}
	if other == nil {
		return false
	}
	return fi.name == other.name
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
