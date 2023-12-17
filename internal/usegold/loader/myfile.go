package loader

import (
	"fmt"
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

func myScanFile(path string) (*MyFile, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return myErrFile(path, fmt.Errorf("file read error (%w)", err))
	}
	return &MyFile{
		myTreeItem: myTreeItem{
			parent: nil,
			name:   path,
		},
		content: contents,
	}, nil
}

// myErrFile returns "fake" markdown file showing an error message.
func myErrFile(path string, err error) (*MyFile, error) {
	return &MyFile{
		myTreeItem: myTreeItem{
			parent: nil,
			name:   path,
		},
		content: []byte(fmt.Sprintf("## Unable to load from %s; %s", path, err.Error())),
	}, err
}
