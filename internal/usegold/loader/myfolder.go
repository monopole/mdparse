package loader

import (
	"log/slog"
)

// MyFolder is a named group of files and folders.
type MyFolder struct {
	myTreeItem
	files []*MyFile
	dirs  []*MyFolder
}

var _ MyTreeItem = &MyFolder{}

func NewFolder(n string) *MyFolder {
	return &MyFolder{myTreeItem: myTreeItem{name: n}}
}

func (fl *MyFolder) Accept(v TreeVisitor) {
	v.VisitFolder(fl)
}

func (fl *MyFolder) AddFileObject(file *MyFile) {
	file.parent = fl
	fl.files = append(fl.files, file)
}

func (fl *MyFolder) AddFolderObject(folder *MyFolder) {
	folder.parent = fl
	fl.dirs = append(fl.dirs, folder)
}

func (fl *MyFolder) IsEmpty() bool {
	return len(fl.files) == 0 && len(fl.dirs) == 0
}

func (fl *MyFolder) Equals(other *MyFolder) bool {
	if fl == nil {
		return other == nil
	}
	if other == nil {
		return false
	}
	if fl.name != other.name {
		return false
	}
	if !EqualFileSlice(fl.files, other.files) {
		return false
	}
	return EqualFolderSlice(fl.dirs, other.dirs)
}

func EqualFileSlice(s1 []*MyFile, s2 []*MyFile) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if !s1[i].Equals(s2[i]) {
			return false
		}
	}
	return true
}

func EqualFolderSlice(s1 []*MyFolder, s2 []*MyFolder) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if !s1[i].Equals(s2[i]) {
			return false
		}
	}
	return true
}

func (fl *MyFolder) loadFileFromFs(name string) error {
	slog.Debug("adding   FILE", "name", name, "parent", fl.FullName())
	for _, fi := range fl.files {
		if fi.Name() == name {
			// Already got it
			return nil
		}
	}
	fi := MyFile{
		myTreeItem: myTreeItem{
			parent: fl,
			name:   name,
		},
	}
	// Do a test read.  inject FsLoader
	//if _, err := fi.Contents(); err != nil {
	//	return err
	//}
	fl.files = append(fl.files, &fi)
	return nil
}
