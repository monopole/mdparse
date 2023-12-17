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

func (fl *MyFolder) AddFile(file *MyFile) {
	file.parent = fl
	fl.files = append(fl.files, file)
}

func (fl *MyFolder) AddFolder(folder *MyFolder) {
	folder.parent = fl
	fl.dirs = append(fl.dirs, folder)
}

func (fl *MyFolder) absorbFile(path string) error {
	slog.Debug("Absorbing   FILE", "path", path, "parent", fl.FullName())
	dir, name := FSplit(path)
	folder := fl
	if dir != "" {
		folder = fl.findOrCreateDir(dir)
	}
	folder.addFile(name)
	return nil
}

func (fl *MyFolder) absorbFolder(path string) error {
	slog.Debug("Absorbing FOLDER", "path", path, "parent", fl.FullName())
	dir, name := FSplit(path)
	folder := fl
	if dir != "" {
		folder = fl.findOrCreateDir(dir)
	}
	folder.addFolder(name)
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
	return folder.addFolder(name)
}

func (fl *MyFolder) addFile(name string) {
	slog.Debug("adding   FILE", "name", name, "parent", fl.FullName())
	fi := MyFile{
		myTreeItem: myTreeItem{
			parent: fl,
			name:   name,
		},
	}
	fl.files = append(fl.files, &fi)
}

func (fl *MyFolder) addFolder(name string) *MyFolder {
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
