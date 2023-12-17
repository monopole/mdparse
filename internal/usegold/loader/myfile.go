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
