package loader

// MyFile is named byte array.
type MyFile struct {
	myTreeItem
	// content []byte
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
