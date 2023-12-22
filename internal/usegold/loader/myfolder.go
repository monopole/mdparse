package loader

import (
	"log/slog"
	"path/filepath"
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

//switch t := v.(type) {
//case string:
//// t is a string
//case int :
//// t is an int
//default:
//// t is some other type that we didn't name.
//}

// LoadFolder loads the folder at the given path.
func LoadFolder(fsl *FsLoader, path string) (*MyFolder, error) {
	// First create a parent tree.
	// If the path is an absolute path, the parent folder is named "/",
	// and has a nil parent.
	// If the path is a relative path, the parent folder is named "",
	//

	var placeholder MyFolder
	dir, name := filepath.Split(path)
	if dir != "" && dir != string(filepath.Separator) {
	}
	folder := &placeholder
	if dir != "" {
		// Create a chain of placeholder parents...
		folder = placeholder.buildParentTree(dir)
	}
	return fsl.loadFolderFromFs(folder, name)
}

func (fl *MyFolder) buildParentTree(path string) *MyFolder {
	dir, name := filepath.Split(path)
	folder := fl
	if dir != "" && dir != string(filepath.Separator) && dir != "." && dir != "./" {
		folder = fl.buildParentTree(dir)
	}
	return folder.insertSubFolder(name)
}

func (fl *MyFolder) insertSubFolder(name string) *MyFolder {
	for _, d := range fl.dirs {
		if d.name == name {
			return d
		}
	}
	dir := MyFolder{
		myTreeItem: myTreeItem{
			parent: fl,
			name:   name,
		},
	}
	fl.dirs = append(fl.dirs, &dir)
	return &dir
}

// AbsorbFileFromDisk assumes the argument is a path to a file.
// The final file will be made available for loading,
// but nothing on the intervening path will be read
// (i.e. no sibling trees).
func (fl *MyFolder) AbsorbFileFromDisk(path string) error {
	slog.Debug("Absorbing   FILE", "path", path, "parent", fl.FullName())
	dir, name := FSplit(path)
	folder := fl
	if dir != "" {
		folder = fl.buildParentTree(dir)
	}
	return folder.loadFileFromFs(name)
}

// AbsorbFolderFromDisk assumes the argument is a path to a folder.
// The final folder and all it's contents will be loaded in,
// but nothing on the intervening path will be read
// (i.e. no sibling trees).
func (fl *MyFolder) AbsorbFolderFromDisk(ts *FsLoader, path string) error {
	slog.Debug("Absorbing FOLDER", "path", path, "parent", fl.FullName())
	dir, name := FSplit(path)
	folder := fl
	if dir != "" {
		folder = fl.buildParentTree(dir)
	}
	child, err := ts.loadFolderFromFs(folder, name)
	if err != nil {
		return err
	}
	if child != nil {
		folder.dirs = append(folder.dirs, child)
	}
	return nil
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
	// Do a test read.
	if _, err := fi.Contents(); err != nil {
		return err
	}
	fl.files = append(fl.files, &fi)
	return nil
}
