package loader

// MyFile is named byte array.
type MyFile struct {
	myTreeItem
	// content []byte ?  yes probably put this back as a reload should
	// just start frrom scratch - the original arg string
	// so its not a reloadd so much as a load afresh, possibly with
	//  a new arg.  also the repo should be loaded into memory
	// and temp space deleted each time.
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
