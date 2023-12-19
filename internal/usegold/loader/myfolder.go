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

// AbsorbFileFromDisk assumes the argument is a path to a file.
// The final file will be made available for loading,
// but nothing on the intervening path will be read
// (i.e. no sibling trees).
func (fl *MyFolder) AbsorbFileFromDisk(path string) error {
	slog.Debug("Absorbing   FILE", "path", path, "parent", fl.FullName())
	dir, name := FSplit(path)
	folder := fl
	if dir != "" {
		folder = fl.findOrCreateDir(dir)
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
		folder = fl.findOrCreateDir(dir)
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

func (fl *MyFolder) findOrCreateDir(path string) *MyFolder {
	dir, name := FSplit(path)
	slog.Debug("findOrCreateDir", "path", path)
	folder := fl
	if dir != "" && dir != string(filepath.Separator) {
		folder = fl.findOrCreateDir(dir)
	}
	for _, item := range folder.dirs {
		if item.name == name {
			slog.Debug("   found folder", "name", name)
			return item
		}
	}
	slog.Debug("   creating folder", "name", name)
	return folder.addPlaceholderFolder(name)
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

func (fl *MyFolder) addPlaceholderFolder(name string) *MyFolder {
	slog.Debug("adding   FOLDER", "name", name, "parent", fl.FullName())
	dir := MyFolder{
		myTreeItem: myTreeItem{
			parent: fl,
			name:   name,
		},
	}
	fl.dirs = append(fl.dirs, &dir)
	return &dir
}
